import warnings
import onnxruntime as ort
from kokoro_onnx import Kokoro
import sounddevice as sd
import os
import time
import numpy as np
from concurrent.futures import ThreadPoolExecutor
import psutil
import socket
from threading import Thread

# Filter numpy warnings about subnormal values
warnings.filterwarnings('ignore', category=UserWarning, module='numpy._core.getlimits')

# =============================================
# System-Aware Optimization Configuration
# =============================================
def get_optimal_thread_count():
    cpu_count = os.cpu_count()
    memory_gb = psutil.virtual_memory().total / (1024 ** 3)

    # Adjust thread counts based on available resources
    if memory_gb >= 16:  # High memory system
        intra_op = max(1, cpu_count - 2)
        inter_op = 2
    else:  # Lower memory system
        intra_op = max(1, cpu_count // 2)
        inter_op = 1
        
    return intra_op, inter_op

INTRA_OP_THREADS, INTER_OP_THREADS = get_optimal_thread_count()

OPTIMIZATION_CONFIG = {
    "intra_op_threads": INTRA_OP_THREADS,
    "inter_op_threads": INTER_OP_THREADS,
    "execution_mode": ort.ExecutionMode.ORT_SEQUENTIAL,
    "optimization_level": ort.GraphOptimizationLevel.ORT_ENABLE_EXTENDED,
    "memory_limit_mb": int(psutil.virtual_memory().available / (1024 * 1024) * 0.7),
    "batch_size": 8192
}

# =============================================
# Environment Optimization
# =============================================
def optimize_environment():
    os.environ["OMP_NUM_THREADS"] = str(OPTIMIZATION_CONFIG["intra_op_threads"])
    os.environ["KMP_BLOCKTIME"] = "1"
    os.environ["KMP_AFFINITY"] = "granularity=fine,compact,1,0"
    os.environ["OMP_WAIT_POLICY"] = "ACTIVE"
    os.environ["OMP_PROC_BIND"] = "TRUE"

# =============================================
# Session Optimization
# =============================================
def create_optimized_session(model_path="./kokoro-tts/voice.onnx"):
    session_options = ort.SessionOptions()
    
    # Thread configuration
    session_options.intra_op_num_threads = OPTIMIZATION_CONFIG["intra_op_threads"]
    session_options.inter_op_num_threads = OPTIMIZATION_CONFIG["inter_op_threads"]
    
    # Performance optimization
    session_options.execution_mode = OPTIMIZATION_CONFIG["execution_mode"]
    session_options.graph_optimization_level = OPTIMIZATION_CONFIG["optimization_level"]
    session_options.enable_cpu_mem_arena = True
    session_options.enable_mem_pattern = True
    session_options.enable_mem_reuse = True
    
    # Disable NCHWc transformation
    session_options.add_session_config_entry("session.disable_nchwc_transformer", "1")
    
    # Memory optimization
    session_options.add_session_config_entry("session.set_denormal_as_zero", "1")
    session_options.add_session_config_entry("session.use_deterministic_compute", "0")
    
    return ort.InferenceSession(
        model_path,
        providers=["CPUExecutionProvider"],
        sess_options=session_options,
        provider_options=[{
            "arena_extend_strategy": "kSameAsRequested",
            "cpu_mem_limits": OPTIMIZATION_CONFIG["memory_limit_mb"] * 1024 * 1024
        }]
    )

# =============================================
# Audio Processing Optimization
# =============================================
def process_audio_batch(samples, sample_rate):
    # Handle potential subnormal values
    with np.errstate(all='ignore'):
        samples = np.float32(samples)
        output = np.zeros_like(samples)
        
        batch_size = OPTIMIZATION_CONFIG["batch_size"]
        for i in range(0, len(samples), batch_size):
            end = min(i + batch_size, len(samples))
            output[i:end] = samples[i:end]
        
        output[np.abs(output) < np.finfo(np.float32).tiny] = 0
        return output

# =============================================
# Client Handling and Processing
# =============================================
def handle_client(conn, kokoro):
    try:
        with conn:
            buffer = ''
            while True:
                data = conn.recv(1024).decode('utf-8')
                if not data:
                    break
                buffer += data
                while '\n' in buffer:
                    text, buffer = buffer.split('\n', 1)
                    text = text.strip()
                    if text:
                        process_text(kokoro, text)
                        conn.sendall(b'OK\n')
    except Exception as e:
        print(f"Connection error: {e}")

def process_text(kokoro, text):
    start_time = time.time()
    
    # Generate audio
    samples, sample_rate = kokoro.create(
        text,
        voice="af_sky",
        speed=1.2,
        lang="en-us",
    )
    
    # Process audio
    with ThreadPoolExecutor(max_workers=2) as executor:
        processed_samples = executor.submit(
            process_audio_batch, samples, sample_rate
        ).result()
    
    generation_time = time.time() - start_time
    print(f"Generation time: {generation_time:.2f} seconds | Text: '{text}'")
    
    # Play audio
    sd.play(processed_samples, sample_rate, blocking=False)
    sd.wait()

# =============================================
# Main Execution
# =============================================
def main():
    try:
        optimize_environment()
        session = create_optimized_session()
        kokoro = Kokoro.from_session(session, "./kokoro-tts/voices.json")

        HOST = 'localhost'
        PORT = 65432
        
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            s.bind((HOST, PORT))
            s.listen()
            print(f"Server listening on {HOST}:{PORT}")
            
            while True:
                conn, addr = s.accept()
                print(f"Connected by {addr}")
                Thread(target=handle_client, args=(conn, kokoro)).start()
                
    except KeyboardInterrupt:
        print("\nShutting down server gracefully...")
    except Exception as e:
        print(f"Critical error: {e}")
    finally:
        if 'session' in locals():
            del session
        if 'kokoro' in locals():
            del kokoro
        print("Cleanup complete")

if __name__ == "__main__":
    main()

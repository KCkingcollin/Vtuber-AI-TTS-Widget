# -*- mode: python ; coding: utf-8 -*-

import os
from PyInstaller.utils.hooks import collect_data_files

block_cipher = None

# Bundle voices.json from the local directory
datas = [
    ('./voices.json', '.'), 
    ('./winvenv/Lib/site-packages/language_tags', 'language_tags'), 
    ('./winvenv/Lib/site-packages/espeakng_loader', 'espeakng_loader')
]

# Automatically collect JSON data files from the language_tags package.
# This ensures that language_tags can find its required JSON files.
datas += collect_data_files('language_tags', subdir='data/json')

a = Analysis(
    ['tts-server.py'],
    pathex=[os.path.abspath('.')],  # Ensure the current directory is in the search path
    binaries=[],
    datas=datas,
    hiddenimports=[],  # Add any additional hidden imports if needed
    hookspath=[],
    runtime_hooks=[],
    excludes=[],
    win_no_prefer_redirects=False,
    win_private_assemblies=False,
    cipher=block_cipher,
)

pyz = PYZ(a.pure, a.zipped_data, cipher=block_cipher)

exe = EXE(
    pyz,
    a.scripts,
    [],
    exclude_binaries=True,
    name='tts-server-win',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    console=False,  # Change to False if you do not want a console window
)

coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=False,
    upx=True,
    name='tts-server-win'
)

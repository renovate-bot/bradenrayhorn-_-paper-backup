declare global {
  interface Window {
    paperBackup: (
      data: Uint8Array,
      passphrase: string,
    ) => { text: string; qr: Uint8Array };

    paperBackupDecodeQR: (passphrase: string, data: Uint8Array) => Uint8Array;
  }
}

export {};

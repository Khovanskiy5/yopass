export function randomInt(min: number, max: number): number {
  const byteArray = new Uint8Array(1);
  window.crypto.getRandomValues(byteArray);

  const range = max - min;
  const maxRange = 256;
  if (byteArray[0] >= Math.floor(maxRange / range) * range) {
    return randomInt(min, max);
  }
  return min + (byteArray[0] % range);
}

export function randomString(): string {
  const charset =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+~';
  const byteArray = new Uint8Array(32);
  window.crypto.getRandomValues(byteArray);
  let text = '';
  for (let i = 0; i < byteArray.length; i++) {
    text += charset[byteArray[i] % charset.length];
  }
  return text;
}

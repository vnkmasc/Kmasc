import LZString from 'lz-string'

export const encodeJSON = (obj: any) => {
  try {
    const str = JSON.stringify(obj)
    return LZString.compressToEncodedURIComponent(str)
  } catch {
    return null
  }
}

export const decodeJSON = (encodedStr: string) => {
  try {
    const str = LZString.decompressFromEncodedURIComponent(encodedStr)
    if (!str) return null
    return JSON.parse(str)
  } catch {
    return null
  }
}

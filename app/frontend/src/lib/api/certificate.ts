import { queryString } from '../utils/common'
import { formatCertificate, formatCertificateVerifyCode, formatCertificateView } from '../utils/format-api'
import apiService from './root'

export const getCertificateList = async (params: any) => {
  const res = await apiService('GET', queryString(['certificates', 'search'], params))
  return {
    ...res,
    data: res.data.map((item: any) => formatCertificate(item))
  }
}

export const createCertificate = async (data: any) => {
  const formattedData = formatCertificate(data, true) as Record<string, string | number | null | undefined>
  const res = await apiService('POST', 'certificates', formattedData)
  return res
}

export const uploadCertificate = async (data: any, name: string) => {
  const res = await apiService('POST', `certificates/upload-pdf?is_degree=false&name=${name}`, data)
  return res
}

export const uploadDegree = async (data: any) => {
  const res = await apiService('POST', 'certificates/upload-pdf?is_degree=true', data)
  return res
}

export const importCertificateExcel = async (data: any) => {
  const res = await apiService('POST', 'certificates/import-excel', data)
  return res
}

export const getCertificateDataById = async (id: string) => {
  const res = await apiService('GET', `certificates/${id}`)

  return formatCertificateView(res.data)
}

export const getCertificateFile = async (id: string) => {
  const res = await apiService('GET', `certificates/file/${id}`, undefined, true, {}, true)
  return res
}

export const getVerifyCodeList = async (params: any) => {
  const res = await apiService('GET', queryString(['verification', 'my-codes'], params))
  return {
    ...res,
    data: res.data.map((item: any) => formatCertificateVerifyCode(item))
  }
}

export const createVerifyCode = async (data: any) => {
  const formattedData = formatCertificateVerifyCode(data, true) as Record<string, string | number | null | undefined>
  const res = await apiService('POST', 'verification/create', formattedData)
  return res
}

export const getCertificateDataStudent = async () => {
  const res = await apiService('GET', 'certificates/my-certificate')

  return formatCertificateView(res.data[0])
}

export const verifyCodeDataforGuest = async (code: string) => {
  const res = await apiService(
    'POST',
    'auth/verification',
    {
      code,
      view_type: 'data'
    },
    false
  )
  return formatCertificateView(res.data)
}

export const verifyCodeFileforGuest = async (code: string) => {
  const res = await apiService(
    'POST',
    'auth/verification',
    {
      code,
      view_type: 'file'
    },
    false,
    {},
    true
  )
  return res
}

export const getCertificatesNameByStudent = async () => {
  const res = await apiService('GET', 'certificates/simple')
  return res.data
}

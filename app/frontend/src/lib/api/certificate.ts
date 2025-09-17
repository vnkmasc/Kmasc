import { queryString } from '../utils/common'
import { formatCertificate, formatCertificateView } from '../utils/format-api'
import apiService from './root'

export const getCertificateList = async (params: any) => {
  const res = await apiService('GET', queryString(['certificates', 'search'], params))
  return {
    ...res,
    data: res.data.map((item: any) => formatCertificate(item, false))
  }
}

export const createCertificate = async (data: any) => {
  const formattedData = formatCertificate(data, false, true)
  const res = await apiService('POST', 'certificates', formattedData)
  return res
}

export const createDegree = async (data: any) => {
  const formattedData = formatCertificate(data, true, true)
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

  return {
    certificate: formatCertificateView(res.data)
  }
}

export const getCertificateFile = async (id: string) => {
  const res = await apiService('GET', `certificates/file/${id}`, undefined, true, {}, true)
  return res
}

export const getCertificateDataStudent = async () => {
  const res = await apiService('GET', 'certificates/my-certificate')

  return formatCertificateView(res.data[0])
}

export const getCertificatesNameByStudent = async () => {
  const res = await apiService('GET', 'certificates/simple')
  return res.data
}

export const pushCertificateIntoBlockchain = async (id: string) => {
  const res = await apiService('POST', `blockchain/push-chain/${id}`)
  return res
}

export const getBlockchainData = async (
  universityId: string,
  facultyId: string,
  certificateType: string,
  course: string,
  certificateId: string
) => {
  const res = await apiService(
    'POST',
    'blockchain/verify-batch-certificates',
    {
      university_id: universityId,
      faculty_id: facultyId,
      certificate_type: certificateType,
      course: course,
      certificate_id: certificateId
    },
    false
  )
  return {
    ...res,
    certificate: formatCertificateView(res.data)
  }
}

export const getBlockchainFile = async (id: string) => {
  const res = await apiService('GET', `certificates/file/${id}`, undefined, true, {}, true)
  return res
}

export const uploadCertificatesBlockchain = async (facultyId: string, certificateType: string, course: string) => {
  const res = await apiService('POST', 'blockchain/push-certificates', {
    faculty_id: facultyId,
    certificate_type: certificateType,
    course: course
  })
  return res
}

export const verifyDegreeDataBlockchain = async (
  universityId: string,
  facultyId: string,
  certificateType: string,
  course: string,
  certificateId: string
) => {
  const res = await apiService(
    'POST',
    'blockchain/verify-batch-certificates',
    {
      university_id: universityId,
      faculty_id: facultyId,
      certificate_type: certificateType,
      course: course,
      certificate_id: certificateId
    },
    false
  )
  return res
}

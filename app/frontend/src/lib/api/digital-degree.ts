import { queryString } from '../utils/common'
import apiService from './root'

export const createDegreeTemplate = async (data: any) => {
  const res = await apiService('POST', 'templates', data)
  return res
}

export const getDegreeTemplateById = async (id: string) => {
  const res = await apiService('GET', `templates/${id}`)

  return {
    ...res.data,
    faculty_id: res.data.facultyId
  }
}

export const updateDegreeTemplate = async (id: string, data: any) => {
  const res = await apiService('PUT', `templates/${id}`, data)
  return res
}
export const searchDegreeTemplateByFaculty = async (facultyId: string) => {
  const res = await apiService('GET', `templates/faculty/${facultyId}`)
  return res
}

export const deleteDegreeTemplate = async (id: string) => {
  const res = await apiService('DELETE', `templates/${id}`)
  return res
}

export const signDegreeTemplateFaculty = async (facultyId: string) => {
  const res = await apiService('POST', `templates/sign/faculty/${facultyId}`)
  return res
}

export const signDegreeTemplateUni = async () => {
  const res = await apiService('POST', `templates/sign/university`)
  return res
}

export const signDegreeTemplateById = async (id: string) => {
  const res = await apiService('POST', `templates/${id}/sign`)
  return res
}

export const searchDigitalDegreeList = async (params: any) => {
  const res = await apiService('GET', queryString(['ediplomas', 'search'], params))
  return res
}

export const issueDownloadDegreeZip = async (facultyId: string, templateId: string, fileName: string) => {
  const blob = await apiService(
    'POST',
    `ediplomas/generate-bulk-zip`,
    {
      faculty_id: facultyId,
      template_id: templateId
    },
    true,
    { Accept: 'application/zip' },
    true
  )

  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  window.URL.revokeObjectURL(url)

  return blob
}

export const uploadDegreeToMinio = async () => {
  const res = await apiService('POST', `ediplomas/upload-local`)
  return res
}

export const getTemplateInterfaces = async () => {
  const res = await apiService('GET', `template-samples`)
  return res
}

export const getTemplateInterfaceById = async (id: string) => {
  const res = await apiService('GET', `template-samples/${id}`)
  return res
}

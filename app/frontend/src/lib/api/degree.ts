import { queryString } from '../utils/common'
import apiService from './root'

export const createDegreeTemplate = async (data: any) => {
  const res = await apiService('POST', 'templates', data)
  return res
}

export const getDegreeTemplateById = async (id: string) => {
  const res = await apiService('GET', `templates/${id}`)
  return res
}

export const updateDegreeTemplate = async (id: string, data: any) => {
  const res = await apiService('PUT', `templates/${id}`, data)
  return res
}
export const searchDegreeTemplateByFaculty = async (facultyId: string) => {
  const res = await apiService('GET', `templates/faculty/${facultyId}`)
  return res
}

export const getDegreeTemplateView = async (id: string) => {
  const res = await apiService('GET', `templates/${id}`)

  return res.data?.html_content
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

export const issueDigitalDegreeFaculty = async (facultyId: string, templateId: string) => {
  const res = await apiService('POST', `ediplomas/generate-bulk`, {
    faculty_id: facultyId,
    template_id: templateId
  })
  return res
}

export const downloadDegreeZip = async (facultyId: string, templateId: string) => {
  const blob = await apiService(
    'POST',
    `ediplomas/generate-bulk-zip`,
    {
      faculty_id: facultyId,
      template_id: templateId
    },
    true,
    { Accept: 'application/zip' },
    true // isBlob = true to get blob response
  )

  // Tự động tải file về
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `degrees-${facultyId}-${templateId}.zip`
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

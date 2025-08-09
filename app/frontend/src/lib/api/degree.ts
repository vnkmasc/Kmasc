import { queryString } from '../utils/common'
import { formatDegreeTemplateFormData } from '../utils/format-api'
import apiService from './root'

export const createDegreeTemplate = async (data: any) => {
  const res = await apiService('POST', 'templates', formatDegreeTemplateFormData(data))
  return res
}

export const getDegreeTemplateById = async (id: string) => {
  const res = await apiService('GET', `templates/${id}`)
  return res
}

export const updateDegreeTemplate = async (id: string, data: any) => {
  const res = await apiService('PUT', `templates/${id}`, formatDegreeTemplateFormData(data, false))
  return res
}
export const searchDegreeTemplateByFaculty = async (facultyId: string) => {
  const res = await apiService('GET', `templates/faculty/${facultyId}`)
  return res
}

export const getDegreeTemplateView = async (id: string) => {
  const res = await apiService('GET', `templates/view/${id}`, undefined, true, { Accept: 'text/html' })

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

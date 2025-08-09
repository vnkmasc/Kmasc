import { formatDegreeTemplateFormData } from '../utils/format-api'
import apiService from './root'

export const createDegreeTemplate = async (data: any) => {
  const res = await apiService('POST', 'templates', formatDegreeTemplateFormData(data))
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

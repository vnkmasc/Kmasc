import { formatFaculty } from '../utils/format-api'
import apiService from './root'

export const getFacultyList = async () => {
  const res = await apiService('GET', 'faculties')
  return res.data.map((item: any) => formatFaculty(item))
}

export const getFacultyById = async (id: string) => {
  const res = await apiService('GET', `faculties/${id}`)
  return formatFaculty(res.data)
}

export const createFaculty = async (data: any) => {
  const res = await apiService('POST', 'faculties', formatFaculty(data, true))
  return res
}

export const updateFaculty = async (id: string, data: any) => {
  const res = await apiService('PUT', `faculties/${id}`, formatFaculty(data, true))
  return res
}

export const deleteFaculty = async (id: string) => {
  const res = await apiService('DELETE', `faculties/${id}`)
  return res
}

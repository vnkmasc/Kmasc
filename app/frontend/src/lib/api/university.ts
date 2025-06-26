import { queryString } from '../utils/common'
import apiService from './root'

export const getUniversityList = async (params?: any) => {
  const res = await apiService('GET', queryString(['universities'], params || {}))
  return res.data
}

export const approveUniversity = async (id: string) => {
  const res = await apiService('POST', `universities/approve-or-reject`, {
    university_id: id,
    action: 'approve'
  })
  return res
}

export const rejectUniversity = async (id: string) => {
  const res = await apiService('POST', `universities/approve-or-reject`, {
    university_id: id,
    action: 'reject'
  })
  return res
}

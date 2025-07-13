import { queryString } from '../utils/common'
import { formatRewardDiscipline } from '../utils/format-api'
import apiService from './root'

export const searchRewardDiscipline = async (params: any) => {
  const res = await apiService('GET', queryString(['reward-disciplines', 'search'], params))
  return {
    ...res,
    data: res.data.map((item: any) => formatRewardDiscipline(item))
  } as any
}

export const createRewardDiscipline = async (data: any) => {
  const res = await apiService('POST', 'reward-disciplines', formatRewardDiscipline(data, true))
  return res
}

export const updateRewardDiscipline = async (id: string, data: any) => {
  const res = await apiService('PUT', `reward-disciplines/${id}`, formatRewardDiscipline(data, true))
  return res
}

export const deleteRewardDiscipline = async (id: string) => {
  const res = await apiService('DELETE', `reward-disciplines/${id}`)
  return res
}

export const getRewardDisciplineById = async (id: string) => {
  const res = await apiService('GET', `reward-disciplines/${id}`)

  return formatRewardDiscipline(res.data)
}

export const importExcelRewardDiscipline = async (isDiscipline: boolean, data: any) => {
  const res = await apiService('POST', `reward-disciplines/import-excel?is_discipline=${isDiscipline}`, data)
  return res
}

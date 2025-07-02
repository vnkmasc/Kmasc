import { queryString } from '../utils/common'
import { formatRewardDiscipline } from '../utils/format-api'
import apiService from './root'

export const searchRewardDiscipline = async (params: any) => {
  const response = await apiService('GET', queryString(['reward-disciplines', 'search'], params))
  return {
    ...response,
    data: response.data.map((item: any) => formatRewardDiscipline(item))
  } as any
}

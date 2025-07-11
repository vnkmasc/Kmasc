import { searchStudentByCode } from '@/lib/api/student'

export const STUDENT_STATUS_OPTIONS = [
  {
    label: 'Đã tốt nghiệp',
    value: 'true'
  },
  {
    label: 'Đang học',
    value: 'false'
  }
]

export const GENDER_SELECT_SETTING = {
  select: {
    groups: [
      {
        label: undefined,
        options: [
          { label: 'Nam', value: 'true' },
          { label: 'Nữ', value: 'false' }
        ]
      }
    ]
  }
}

export const CERTIFICATE_TYPE_OPTIONS = [
  { value: '1', label: 'Cử nhân' },
  { value: '2', label: 'Kỹ sư' },
  { value: '3', label: 'Thạc sĩ' },
  { value: '4', label: 'Tiến sĩ' }
]

export const REWARD_DISCIPLINE_TYPE_SETTING = {
  select: {
    groups: [
      {
        label: undefined,
        options: [
          { label: 'Khen thưởng', value: 'false' },
          { label: 'Kỷ luật', value: 'true' }
        ]
      }
    ]
  }
}

export const STUDENT_CODE_SEARCH_SETTING = {
  querySelect: {
    queryFn: (keyword: string) => searchStudentByCode(keyword)
  }
}
export const PAGE_SIZE = 10

export const LEVEL_DISCIPLINE = {
  1: 'Khiển trách',
  2: 'Cảnh cáo',
  3: 'Đình chỉ tạm thời',
  4: 'Buộc thôi học'
}

export const REWARD_DISCIPLINE_LEVEL_SETTING = {
  select: {
    groups: [
      {
        label: undefined,
        options: Object.entries(LEVEL_DISCIPLINE).map(([key, value]) => ({ label: value, value: key }))
      }
    ]
  }
}

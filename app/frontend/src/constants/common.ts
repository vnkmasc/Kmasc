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
  { value: 'Cử nhân', label: 'Cử nhân' },
  { value: 'Kỹ sư', label: 'Kỹ sư' },
  { value: 'Thạc sĩ', label: 'Thạc sĩ' },
  { value: 'Tiến sĩ', label: 'Tiến sĩ' }
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

export const DEGREE_TEMPLATE_STATUS = {
  PENDING: {
    variant: 'outline',
    label: 'Chưa ký'
  },
  SIGNED_BY_UNI: {
    variant: 'secondary',
    label: 'Đã ký bởi trường đại học'
  },
  SIGNED_BY_MINEDU: {
    variant: 'default',
    label: 'Đã ký bởi TĐF và Bộ GD'
  }
}

export const GRADUATION_RANK_OPTIONS = [
  { value: 'Giỏi', label: 'Giỏi' },
  { value: 'Khá', label: 'Khá' },
  { value: 'Trung bình', label: 'Trung bình' },
  { value: 'Yếu', label: 'Yếu' }
]

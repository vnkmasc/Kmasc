import { HTMLInputTypeAttribute } from 'react'
import { Control } from 'react-hook-form'
import { z } from 'zod'

export interface SessionPayload {
  access_token: string
  role: string
}

export interface OptionType {
  value: string
  label: string
}

interface SelectGroup {
  label: string | undefined
  options: OptionType[]
}

export interface CustomFormItem {
  type: 'input' | 'select' | 'query_select' | 'search_select' | 'textarea' | 'file'
  control: Control<any, any, any> | undefined
  name: string
  label?: string | ''
  placeholder?: string
  description?: string | undefined
  disabled?: boolean | false

  setting?: {
    input?: {
      type?: HTMLInputTypeAttribute
    }
    select?: {
      groups?: SelectGroup[]
    }
    querySelect?: {
      // eslint-disable-next-line no-unused-vars
      queryFn: (keyword: string) => Promise<any>
    }
    file?: {
      accept?: string
    }
    textarea?: {
      rows?: number
    }
  }
}

export interface CustomZodFormItem extends Omit<CustomFormItem, 'control'> {
  validator?: z.ZodType
  defaultValue?: any
}

export type CertificateType = {
  id: string
  universityName: string
  universityCode: string
  facultyName: string
  facultyCode: string
  studentCode: string
  studentName: string
  name: string
  date: string
  certificateType?: string
  signed: boolean
  serialNumber?: string
  regNo?: string
  isDegree?: boolean
}

export type DegreeTemplateType = {
  name: string
  description?: string
  faculty_id: string
  html_content: string
}

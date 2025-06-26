import { ExternalToast, toast } from 'sonner'

// eslint-disable-next-line no-unused-vars
export function debounce<T extends (...args: any[]) => void>(func: T, wait: number): (...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout> | null = null

  return (...args: Parameters<T>) => {
    if (timeout !== null) {
      clearTimeout(timeout)
    }
    timeout = setTimeout(() => {
      func(...args)
    }, wait)
  }
}

export const clearFalsyValueObject = (obj: Record<string, any>) => {
  return Object.fromEntries(Object.entries(obj).filter((entry) => entry[1] !== null && entry[1] !== undefined))
}

export const queryString = (slashParams: (string | number)[], params?: any) => {
  const filteredParams = params ? clearFalsyValueObject(params) : null
  const queryString = filteredParams ? new URLSearchParams(filteredParams as Record<string, string>).toString() : null
  return `${slashParams.join('/')}${queryString ? '?' + queryString : ''}`
}

export const showNotification = (
  type: 'success' | 'error' | 'info' | 'warning' | 'message',
  description: string,
  setting?: ExternalToast
) => {
  return toast[type]('Thông báo', {
    description:
      description ||
      {
        success: 'Thao tác thành công',
        error: 'Thao tác thất bại',
        info: 'Thông tin',
        warning: 'Cảnh báo',
        message: 'Tin nhắn'
      }[type],
    classNames: {
      success: '[&_svg]:!text-green-500',
      error: '[&_svg]:!text-red-500',
      info: '[&_svg]:!text-blue-500',
      warning: '[&_svg]:!text-yellow-500'
    },
    ...setting
  })
}

export const showMessage = (description: string, setting?: ExternalToast) => {
  return toast(description, {
    position: 'top-center',
    ...setting
  })
}

export const formatResponseImportExcel = (
  data: any
): {
  success: number[]
  error: { title: string; row: number[] }[]
} => {
  return {
    success: data.data.success?.map((item: any) => item.row),
    error: Object.values(
      data.data.error?.reduce((acc: Record<string, { title: string; row: number[] }>, item: any) => {
        if (!acc[item.error]) {
          acc[item.error] = { title: item.error, row: [] }
        }
        acc[item.error].row.push(item.row)
        return acc
      }, {})
    )
  }
}

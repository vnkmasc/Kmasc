'use client'

import { Dispatch, SetStateAction, useEffect, useState } from 'react'
import DetailDialog from '../role/education-admin/detail-dialog'
import { validateNoEmpty } from '@/lib/utils/validators'
import { CustomZodFormItem } from '@/types/common'
import { getDataStorage, saveDataStorage } from '@/lib/utils/handle-storage'
import { showNotification } from '@/lib/utils/common'

interface Props {
  open: boolean
  onOpenChange: Dispatch<SetStateAction<boolean>>
  role: string
}

const SignSetting: React.FC<Props> = (props) => {
  const [data, setData] = useState<any>(null)
  const handleSubmit = (values: any) => {
    saveDataStorage('setting', values)
    showNotification('success', 'Cập nhật cài đặt ký số thành công')
    props.onOpenChange(false)
  }

  useEffect(() => {
    const data = getDataStorage('setting')
    if (data) {
      setData(data)
    }
  }, [])

  const formItems: CustomZodFormItem[] =
    props.role === 'university_admin'
      ? [
          {
            type: 'input',
            label: 'Link server ký số',
            name: 'signService',
            placeholder: 'Nhập link server ký số',
            validator: validateNoEmpty('Link server ký số'),
            description: 'Link server để tiến hành xác minh chữ ký số.'
          },
          {
            type: 'input',
            label: 'Đường dẫn ứng dụng ký PDF',
            name: 'pdfSignLocation',
            placeholder: 'Nhập đường dẫn ứng dụng chữ ký PDF',
            validator: validateNoEmpty('Đường dẫn ứng dụng ký PDF'),
            description: 'Đường dẫn ứng dụng để tiến hành ký PDF.'
          }
        ]
      : [
          {
            type: 'input',
            label: 'Link server ký số',
            name: 'signService',
            placeholder: 'Nhập link server ký số',
            validator: validateNoEmpty('Link server ký số'),
            description: 'Link server để tiến hành xác minh chữ ký số.'
          }
        ]

  return (
    <DetailDialog
      title='Cài đặt cấu hình ký số'
      handleSubmit={handleSubmit}
      handleClose={() => props.onOpenChange(false)}
      data={data}
      mode={props.open ? 'update' : undefined}
      items={formItems}
    />
  )
}

export default SignSetting

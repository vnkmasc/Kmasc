'use client'

import { Dispatch } from 'react'
import { validateNoEmpty } from '@/lib/utils/validators'
import useSWRMutation from 'swr/mutation'
import { changePassword } from '@/lib/api/auth'

import DetailDialog from '../role/education-admin/detail-dialog'
import { showNotification } from '@/lib/utils/common'

interface Props {
  open: boolean
  setOpen: Dispatch<React.SetStateAction<boolean>>
}

const ChangePassDialog: React.FC<Props> = (props) => {
  const mutateChangePassword = useSWRMutation(
    'change-password',
    (_, { arg }: { arg: { oldPassword: string; newPassword: string } }) =>
      changePassword(arg.oldPassword, arg.newPassword),
    {
      onSuccess: () => {
        showNotification('success', 'Thay đổi mật khẩu thành công')
        props.setOpen(false)
      },
      onError: (error) => {
        showNotification('error', error.message || 'Thay đổi mật khẩu thất bại')
      }
    }
  )

  const handleSubmit = (data: { oldPassword: string; newPassword: string }) => {
    mutateChangePassword.trigger(data)
  }

  return (
    <DetailDialog
      data={[]}
      mode={props.open ? 'update' : undefined}
      handleClose={() => props.setOpen(false)}
      handleSubmit={handleSubmit}
      items={[
        {
          type: 'input',
          name: 'oldPassword',
          label: 'Mật khẩu cũ',
          placeholder: 'Nhập mật khẩu cũ',
          setting: { input: { type: 'password' } },
          validator: validateNoEmpty('Mật khẩu cũ')
        },
        {
          type: 'input',
          name: 'newPassword',
          label: 'Mật khẩu mới',
          placeholder: 'Nhập mật khẩu mới',
          setting: { input: { type: 'password' } },
          validator: validateNoEmpty('Mật khẩu mới')
        }
      ]}
      title='Thay đổi mật khẩu'
    />
  )
}

export default ChangePassDialog

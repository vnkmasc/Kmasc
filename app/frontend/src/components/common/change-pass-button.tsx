'use client'
import { RotateCcwKeyIcon } from 'lucide-react'
import { Button } from '../ui/button'
import { useState } from 'react'
import { validateNoEmpty } from '@/lib/utils/validators'
import useSWRMutation from 'swr/mutation'
import { changePassword } from '@/lib/api/auth'

import DetailDialog from '../role/education-admin/detail-dialog'
import { showNotification } from '@/lib/utils/common'

interface Props {
  className?: string
}

const ChangePassButton: React.FC<Props> = (props) => {
  const [open, setOpen] = useState(false)

  const mutateChangePassword = useSWRMutation(
    'change-password',
    (_, { arg }: { arg: { oldPassword: string; newPassword: string } }) =>
      changePassword(arg.oldPassword, arg.newPassword),
    {
      onSuccess: () => {
        showNotification('success', 'Thay đổi mật khẩu thành công')
        setOpen(false)
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
    <>
      <Button variant={'secondary'} size={'icon'} onClick={() => setOpen(true)} className={props.className}>
        <RotateCcwKeyIcon />
      </Button>

      <DetailDialog
        data={[]}
        mode={open ? 'update' : undefined}
        handleClose={() => setOpen(false)}
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
    </>
  )
}

export default ChangePassButton

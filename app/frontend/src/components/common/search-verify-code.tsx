'use client'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'

import { useState } from 'react'
import useSWRMutation from 'swr/mutation'

import { verifyCodeDataforGuest } from '@/lib/api/certificate'
import { showNotification } from '@/lib/utils/common'
import { PackageSearch } from 'lucide-react'

const SearchVerifyCode = () => {
  const [verifyCode, setVerifyCode] = useState('')
  const mutateDataVerifyCode = useSWRMutation(
    'verify-code',
    (_, { arg }: { arg: string }) => verifyCodeDataforGuest(arg),
    {
      onError(error) {
        showNotification('error', error.message || 'Lỗi khi xác thực mã xác minh')
      }
    }
  )

  return (
    <>
      <Card className='mt-6 w-full max-w-[600px] md:mt-10'>
        <CardHeader>
          <CardTitle className='px-3 text-center md:px-6'>
            <h3>Nhập mã xác minh để xem chứng chỉ</h3>
          </CardTitle>
        </CardHeader>
        <CardContent className='px-3 md:px-6'>
          <div className='flex items-center gap-2'>
            <Input placeholder='Nhập số chứng chỉ' value={verifyCode} onChange={(e) => setVerifyCode(e.target.value)} />
            <Button
              isLoading={mutateDataVerifyCode.isMutating}
              onClick={() => mutateDataVerifyCode.trigger(verifyCode)}
            >
              <PackageSearch /> <span className='hidden md:block'>Xác thực</span>
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* {mutateDataVerifyCode.data && (
        <Card className='mt-6 w-full max-w-[800px]'>
          <CardHeader>
            <CardTitle>
              <div className='flex items-center justify-between'>
                <span>Thông tin văn bằng - chứng chỉ</span>
                <CertificateBlankButton isIcon={md ? false : true} action={() => verifyCodeFileforGuest(verifyCode)} />
              </div>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <CertificateView
              data={mutateDataVerifyCode.data as CertificateType}
              className='grid-cols-1 md:grid-cols-2'
            />
          </CardContent>
        </Card>
      )} */}
    </>
  )
}

export default SearchVerifyCode

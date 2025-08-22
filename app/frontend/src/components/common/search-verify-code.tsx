'use client'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useEffect, useMemo, useState } from 'react'
import { PackageSearch } from 'lucide-react'
import CertificateView from './certificate-view'
import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { decodeJSON } from '@/lib/utils/lz-string'
import DigitalDegreeView from './digital-degree-view'

const SearchVerifyCode = () => {
  const [verifyCode, setVerifyCode] = useState('')
  const [blockchainVerifyCode, setBlockchainVerifyCode] = useState('')
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()
  const decodedVerifyCode = useMemo(() => {
    if (!blockchainVerifyCode) return null
    return decodeJSON(blockchainVerifyCode)
  }, [blockchainVerifyCode])

  useEffect(() => {
    if (blockchainVerifyCode) {
      const currentCode = searchParams.get('code')
      if (currentCode !== blockchainVerifyCode) {
        router.push(`${pathname}?code=${blockchainVerifyCode}`)
      }
    }
  }, [blockchainVerifyCode, pathname, router, searchParams])

  useEffect(() => {
    const code = searchParams.get('code')
    if (code) {
      setBlockchainVerifyCode(code)
      setVerifyCode(code)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <>
      <Card className='my-6 w-full max-w-[600px] md:mt-10'>
        <CardHeader>
          <CardTitle className='px-3 text-center md:px-6'>
            <h3>Nhập mã xác minh để xem chứng chỉ</h3>
          </CardTitle>
        </CardHeader>
        <CardContent className='px-3 md:px-6'>
          <div className='flex items-center gap-2'>
            <Input placeholder='Nhập số chứng chỉ' value={verifyCode} onChange={(e) => setVerifyCode(e.target.value)} />
            <Button onClick={() => setBlockchainVerifyCode(verifyCode)}>
              <PackageSearch /> <span className='hidden md:block'>Xác thực</span>
            </Button>
          </div>
        </CardContent>
      </Card>
      <div className='w-full'>
        {!blockchainVerifyCode ? null : decodedVerifyCode ? (
          <DigitalDegreeView
            id={decodedVerifyCode?.ediploma_id}
            isBlockchain={true}
            universityCode={decodedVerifyCode?.university_code}
            universityId={decodedVerifyCode?.university_id}
            facultyId={decodedVerifyCode?.faculty_id}
            certificateType={decodedVerifyCode?.certificate_type}
            course={decodedVerifyCode?.course}
          />
        ) : (
          <CertificateView id={blockchainVerifyCode} isBlockchain={true} />
        )}
      </div>
    </>
  )
}

export default SearchVerifyCode

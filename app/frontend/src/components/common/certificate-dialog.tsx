'use client'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '../ui/dialog'
import { Button } from '../ui/button'
import { PlusIcon } from 'lucide-react'
import { useState } from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select'
import { Label } from '../ui/label'

const CertificateCreateButton: React.FC = () => {
  const [isDegree, setIsDegree] = useState(false)
  console.log('🚀 ~ isDegree:', isDegree)

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>
          <PlusIcon />
          Tạo mới{' '}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Tạo mới văn bằng/chứng chỉ</DialogTitle>
        </DialogHeader>
        <Label>Chọn loại</Label>
        <Select defaultValue='degree' onValueChange={(value) => setIsDegree(value === 'degree')}>
          <SelectTrigger>
            <SelectValue placeholder='Chọn loại văn bằng/chứng chỉ' />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value='degree'>Văn bằng</SelectItem>
            <SelectItem value='certificate'>Chứng chỉ</SelectItem>
          </SelectContent>
        </Select>
        {/* <Form {...form}>
          <form onSubmit={form.handleSubmit(props.handleSubmit)} className='space-y-4'>
            {props.items.map((prop, index) => (
              <CustomFormItem {...prop} control={form.control} key={index} />
            ))}
            <DialogFooter>
              <DialogClose asChild>
                <Button variant='outline' type='button'>
                  Hủy bỏ
                </Button>
              </DialogClose>
              <Button type='submit'>{localMode === 'create' ? 'Tạo mới' : 'Cập nhật'}</Button>
            </DialogFooter>
          </form>
        </Form> */}
      </DialogContent>
    </Dialog>
  )
}

export default CertificateCreateButton

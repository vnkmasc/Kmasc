'use client'

import { Button } from '../ui/button'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogTitle,
  AlertDialogTrigger
} from '../ui/alert-dialog'
import { LogOutIcon } from 'lucide-react'
import { signOut } from '@/lib/auth/auth'
const SignOutButton = () => {
  const handleSignOut = () => {
    signOut()
  }
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button variant='destructive' size={'icon'}>
          <LogOutIcon />
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent className='rounded-lg'>
        <AlertDialogTitle>Đăng xuất</AlertDialogTitle>
        <AlertDialogDescription>Bạn có chắc chắn muốn đăng xuất không?</AlertDialogDescription>
        <AlertDialogFooter>
          <AlertDialogCancel>Hủy bỏ</AlertDialogCancel>
          <AlertDialogAction onClick={handleSignOut}>Đăng xuất</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default SignOutButton

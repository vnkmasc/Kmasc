'use server'

import { redirect } from 'next/navigation'
import { createSession, deleteSession } from './session'
import apiService from '../api/root'

export async function signOut() {
  await deleteSession()
  redirect('/')
}

export async function signIn(payload: { email: string; password: string }) {
  try {
    const data = await apiService('POST', 'auth/login', payload, false)

    if (!data.token) {
      return false
    }

    await createSession({
      access_token: data.token,
      role: data.role
    })

    return true
  } catch (error) {
    console.error('Error sign in', error)
    return false
  }
}

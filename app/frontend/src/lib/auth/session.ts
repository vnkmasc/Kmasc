'use server'

import { SignJWT, jwtVerify } from 'jose'
import { SessionPayload } from '@/types/common'
import { cookies } from 'next/headers'

const sessionKey = process.env.SESSION_SECRET
const encodedKey = new TextEncoder().encode(sessionKey)

export async function encrypt(payload: any) {
  return new SignJWT(payload)
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime('1d')
    .sign(encodedKey)
}

export async function decrypt(session: string | undefined = '') {
  try {
    const { payload } = await jwtVerify(session, encodedKey, {
      algorithms: ['HS256']
    })

    return payload
  } catch (error) {
    console.error(error)
  }
}

export async function createSession(payload: SessionPayload) {
  const expiresAt = new Date(Date.now() + 1000 * 60 * 60 * 24)
  const session = await encrypt(payload)

  const cookieStore = await cookies()
  cookieStore.set('session', session, {
    httpOnly: true,
    expires: expiresAt,
    sameSite: 'lax',
    path: '/'
  })
}

export async function deleteSession() {
  const cookieStore = await cookies()
  cookieStore.delete('session')
}

export async function getSession() {
  const cookieStore = await cookies()
  const session = cookieStore.get('session')

  return session ? decrypt(session.value) : null
}

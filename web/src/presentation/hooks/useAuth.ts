"use client"

import { useSyncExternalStore } from "react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { authRepository } from "@/src/data/repositories/AuthRepository"
import { LoginUseCase }       from "@/src/domain/usecases/LoginUseCase"
import { RegisterUseCase }    from "@/src/domain/usecases/RegisterUseCase"
import { GetMeUseCase }       from "@/src/domain/usecases/GetMeUseCase"
import { VerifyEmailUseCase } from "@/src/domain/usecases/VerifyEmailUseCase"
import { AuthService }        from "@/src/domain/services/AuthService"
import type { LoginInput, RegisterInput } from "@/src/domain/entities/User"

const loginUseCase       = new LoginUseCase(authRepository)
const registerUseCase    = new RegisterUseCase(authRepository)
const getMeUseCase       = new GetMeUseCase(authRepository)
const verifyEmailUseCase = new VerifyEmailUseCase(authRepository)

export function useAuthToken() {
  return useSyncExternalStore(
    AuthService.subscribe,
    () => AuthService.getToken(),
    () => null,
  )
}

export function useMe() {
  const token = useAuthToken()

  const query = useQuery({
    queryKey: ["me"],
    queryFn: () => getMeUseCase.execute(),
    enabled: Boolean(token),
    retry: false,
    staleTime: 5 * 60_000,
  })

  return {
    ...query,
    data: token ? query.data : undefined,
  }
}

export function useLogin() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: LoginInput) => loginUseCase.execute(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  })
}

export function useRegister() {
  return useMutation({
    mutationFn: (input: RegisterInput) => registerUseCase.execute(input),
  })
}

export function useVerifyEmail() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (token: string) => verifyEmailUseCase.execute(token),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  })
}

export function useUpdateProfile() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (input: { display_name?: string; avatar_url?: string }) => authRepository.updateProfile(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  })
}

export function useDeleteMe() {
  return useMutation({
    mutationFn: () => authRepository.deleteMe(),
  })
}

export function useSetUsername() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (username: string) => authRepository.setUsername(username),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["me"] }),
  })
}

export function useLogout() {
  const qc = useQueryClient()
  return () => {
    AuthService.clearToken()
    qc.clear()
  }
}

export function useProviders() {
  return useQuery({
    queryKey: ["auth-providers"],
    queryFn: () => authRepository.providers(),
    staleTime: Infinity,
  })
}

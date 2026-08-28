import { sendRequest } from '@/features/send-request/api'
import { useAsync } from '@/hooks'

/** `run(input)` bhejta hai; `data`/`error`/`isLoading` response state deta hai. */
export function useSendRequest() {
  return useAsync(sendRequest)
}

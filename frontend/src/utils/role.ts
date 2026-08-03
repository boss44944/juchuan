/**
 * 角色判断：区分「服务端」与「客户端」访问。
 *
 * - 服务端 = 在运行菊传的电脑本机访问（localhost / 127.0.0.1）
 * - 客户端 = 局域网内其他设备（手机）通过 IP 访问
 */
export function isServerAccess(): boolean {
  if (typeof window === 'undefined') return false
  const h = window.location.hostname
  return (
    h === 'localhost' ||
    h === '127.0.0.1' ||
    h === '::1' ||
    h === '[::1]' ||
    h === ''
  )
}

export function currentRole(): 'server' | 'client' {
  return isServerAccess() ? 'server' : 'client'
}

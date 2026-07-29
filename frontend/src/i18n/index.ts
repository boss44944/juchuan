import { createI18n } from 'vue-i18n'

const messages = {
  'zh-CN': {
    menu: {
      home: '首页',
      history: '历史',
      devices: '设备',
      config: '配置',
      exit: '退出程序'
    },
    error: {
      AUTH_PASSWORD_INVALID: '密码错误',
      FILE_TOO_LARGE: '文件超过限制'
    }
  },
  'en-US': {
    menu: {
      home: 'Home',
      history: 'History',
      devices: 'Devices',
      config: 'Settings',
      exit: 'Exit'
    },
    error: {
      AUTH_PASSWORD_INVALID: 'Invalid password',
      FILE_TOO_LARGE: 'File is too large'
    }
  },
  'ja-JP': {
    menu: {
      home: 'ホーム',
      history: '履歴',
      devices: 'デバイス',
      config: '設定',
      exit: '終了'
    },
    error: {
      AUTH_PASSWORD_INVALID: 'パスワードが正しくありません',
      FILE_TOO_LARGE: 'ファイルサイズが大きすぎます'
    }
  }
}

export default createI18n({
  legacy: false,
  locale: 'zh-CN',
  fallbackLocale: 'en-US',
  messages
})

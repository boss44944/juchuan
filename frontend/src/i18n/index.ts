import { createI18n } from 'vue-i18n'

const localeFromStorage = localStorage.getItem('juchuan_locale') || 'zh-CN'

const messages = {
  'zh-CN': {
    menu: {
      home: '首页',
      history: '历史',
      devices: '设备',
      messages: '消息',
      send: '发送',
      config: '配置',
      exit: '退出程序',
      logout: '登出'
    },
    login: {
      title: '登录',
      deviceName: '设备名称',
      password: '访问密码',
      submit: '进入',
      qrTip: '手机扫码即可访问',
      qrAlt: '登录二维码',
      deviceNameRequired: '请输入设备名称',
      failed: '登录失败，请检查密码'
    },
    devices: {
      columns: {
        name: '设备名称',
        platform: '平台',
        status: '状态',
        actions: '操作'
      },
      status: {
        online: '在线',
        offline: '离线'
      },
      actions: {
        rename: '重命名',
        remove: '删除'
      },
      dialog: {
        renameTitle: '重命名设备',
        renameInput: '输入新的设备名',
        removeTitle: '删除设备',
        removeConfirm: '确认删除这个设备吗？'
      },
      toast: {
        renamed: '设备名称已更新',
        removed: '设备已删除'
      },
      entryTitle: '访问地址',
      copyAddress: '复制地址',
      copySuccess: '地址已复制',
      qrAlt: '访问二维码'
    },
    send: {
      title: '发送',
      selectTargets: '选择设备',
      inputText: '输入文字',
      sendText: '发送文字',
      selectFile: '选择文件',
      sendFile: '发送文件',
      toast: {
        textRequired: '请输入要发送的文字',
        targetRequired: '请选择目标设备',
        fileRequired: '请先选择文件',
        textSent: '文字已发送',
        fileSent: '文件已发送'
      }
    },
    messagesPage: {
      filters: {
        typeAll: '全部类型',
        typeText: '文本',
        typeFile: '文件',
        statusAll: '全部状态',
        deviceId: '设备ID筛选',
        sender: '发送方',
        target: '接收方',
        apply: '筛选',
        reset: '重置',
        mineOnly: '仅看与我相关'
      },
      actions: {
        batchRead: '批量标记已读',
        batchRetry: '批量重发',
        markRead: '标记已读',
        retry: '重发'
      },
      table: {
        time: '时间',
        type: '类型',
        content: '内容',
        download: '下载文件',
        sender: '发送方',
        target: '接收方',
        status: '状态',
        operation: '操作'
      },
      status: {
        created: '待投递',
        delivered: '已送达',
        read: '已读'
      },
      toast: {
        markReadSuccess: '已标记为已读',
        markReadFailed: '标记已读失败',
        retrySuccess: '重发成功',
        retryFailed: '重发失败，请稍后重试',
        loadFailed: '加载消息失败',
        noneMarkable: '没有可标记的消息',
        batchReadSuccess: '已标记 {count} 条消息',
        batchReadFailed: '批量标记失败',
        selectFirst: '请先选择消息',
        batchRetrySuccess: '已重发 {count} 条消息',
        batchRetryFailed: '批量重发失败',
        batchPartial: '操作完成：成功 {success} 条，失败 {failed} 条'
      }
    },
    configPage: {
      labels: {
        port: '服务端口',
        autoOpen: '自动打开',
        password: '访问密码',
        language: '语言'
      },
      languages: {
        zhCN: '中文',
        enUS: 'English',
        jaJP: '日本語'
      },
      save: '保存',
      saved: '配置已保存'
    },
    error: {
      AUTH_PASSWORD_INVALID: '密码错误',
      AUTH_REQUIRED: '请先登录',
      FILE_TOO_LARGE: '文件超过限制',
      INVALID_REQUEST: '请求参数不正确',
      DEVICE_INFO_REQUIRED: '设备信息不完整',
      DEVICE_NAME_EXISTS: '设备名称已存在',
      DEVICE_NOT_FOUND: '设备不存在',
      DEVICE_NAME_REQUIRED: '设备名称不能为空',
      INVALID_STATUS: '状态参数不合法',
      STATUS_UPDATE_FAILED: '状态更新失败',
      INVALID_FILE: '文件格式无效',
      FILE_REQUIRED: '请先选择文件',
      FILE_SAVE_FAILED: '文件保存失败',
      FILE_RECORD_FAILED: '文件记录保存失败',
      FILE_NOT_FOUND: '文件不存在',
      MESSAGE_REQUIRED: '消息内容或目标设备不能为空',
      MESSAGE_SAVE_FAILED: '消息保存失败',
      MESSAGE_LIST_FAILED: '消息列表加载失败',
      NOT_FOUND: '请求的资源不存在',
      CONFLICT: '请求冲突，请重试',
      SERVER_ERROR: '服务内部错误，请稍后重试',
      UNKNOWN: '请求失败，请稍后重试'
    }
  },
  'en-US': {
    menu: {
      home: 'Home',
      history: 'History',
      devices: 'Devices',
      messages: 'Messages',
      send: 'Send',
      config: 'Settings',
      exit: 'Exit',
      logout: 'Logout'
    },
    login: {
      title: 'Sign In',
      deviceName: 'Device Name',
      password: 'Access Password',
      submit: 'Enter',
      qrTip: 'Scan with mobile device to access',
      qrAlt: 'Login QR code',
      deviceNameRequired: 'Please enter a device name',
      failed: 'Login failed, please check password'
    },
    devices: {
      columns: {
        name: 'Device Name',
        platform: 'Platform',
        status: 'Status',
        actions: 'Actions'
      },
      status: {
        online: 'Online',
        offline: 'Offline'
      },
      actions: {
        rename: 'Rename',
        remove: 'Delete'
      },
      dialog: {
        renameTitle: 'Rename Device',
        renameInput: 'Enter new device name',
        removeTitle: 'Delete Device',
        removeConfirm: 'Are you sure you want to delete this device?'
      },
      toast: {
        renamed: 'Device name updated',
        removed: 'Device deleted'
      },
      entryTitle: 'Access URL',
      copyAddress: 'Copy URL',
      copySuccess: 'URL copied',
      qrAlt: 'Access QR code'
    },
    send: {
      title: 'Send',
      selectTargets: 'Select Devices',
      inputText: 'Enter text',
      sendText: 'Send Text',
      selectFile: 'Choose File',
      sendFile: 'Send File',
      toast: {
        textRequired: 'Please enter text to send',
        targetRequired: 'Please select target devices',
        fileRequired: 'Please choose a file first',
        textSent: 'Text sent',
        fileSent: 'File sent'
      }
    },
    messagesPage: {
      filters: {
        typeAll: 'All Types',
        typeText: 'Text',
        typeFile: 'File',
        statusAll: 'All Status',
        deviceId: 'Filter by Device ID',
        sender: 'Sender',
        target: 'Receiver',
        apply: 'Filter',
        reset: 'Reset',
        mineOnly: 'Only Related to Me'
      },
      actions: {
        batchRead: 'Mark Selected as Read',
        batchRetry: 'Retry Selected',
        markRead: 'Mark Read',
        retry: 'Retry'
      },
      table: {
        time: 'Time',
        type: 'Type',
        content: 'Content',
        download: 'Download',
        sender: 'Sender',
        target: 'Receiver',
        status: 'Status',
        operation: 'Operation'
      },
      status: {
        created: 'Pending',
        delivered: 'Delivered',
        read: 'Read'
      },
      toast: {
        markReadSuccess: 'Marked as read',
        markReadFailed: 'Failed to mark as read',
        retrySuccess: 'Resent successfully',
        retryFailed: 'Resend failed, please try again later',
        loadFailed: 'Failed to load messages',
        noneMarkable: 'No messages can be marked',
        batchReadSuccess: 'Marked {count} messages',
        batchReadFailed: 'Batch mark failed',
        selectFirst: 'Please select messages first',
        batchRetrySuccess: 'Resent {count} messages',
        batchRetryFailed: 'Batch resend failed',
        batchPartial: 'Done: {success} succeeded, {failed} failed'
      }
    },
    configPage: {
      labels: {
        port: 'Server Port',
        autoOpen: 'Open Browser Automatically',
        password: 'Access Password',
        language: 'Language'
      },
      languages: {
        zhCN: 'Chinese',
        enUS: 'English',
        jaJP: 'Japanese'
      },
      save: 'Save',
      saved: 'Configuration saved'
    },
    error: {
      AUTH_PASSWORD_INVALID: 'Invalid password',
      AUTH_REQUIRED: 'Please sign in first',
      FILE_TOO_LARGE: 'File is too large',
      INVALID_REQUEST: 'Invalid request parameters',
      DEVICE_INFO_REQUIRED: 'Incomplete device information',
      DEVICE_NAME_EXISTS: 'Device name already exists',
      DEVICE_NOT_FOUND: 'Device not found',
      DEVICE_NAME_REQUIRED: 'Device name is required',
      INVALID_STATUS: 'Invalid status parameter',
      STATUS_UPDATE_FAILED: 'Failed to update status',
      INVALID_FILE: 'Invalid file format',
      FILE_REQUIRED: 'Please select a file first',
      FILE_SAVE_FAILED: 'Failed to save file',
      FILE_RECORD_FAILED: 'Failed to persist file metadata',
      FILE_NOT_FOUND: 'File not found',
      MESSAGE_REQUIRED: 'Message content or targets are required',
      MESSAGE_SAVE_FAILED: 'Failed to save message',
      MESSAGE_LIST_FAILED: 'Failed to load message list',
      NOT_FOUND: 'Requested resource was not found',
      CONFLICT: 'Request conflict, please retry',
      SERVER_ERROR: 'Internal server error, please retry later',
      UNKNOWN: 'Request failed, please try again later'
    }
  },
  'ja-JP': {
    menu: {
      home: 'ホーム',
      history: '履歴',
      devices: 'デバイス',
      messages: 'メッセージ',
      send: '送信',
      config: '設定',
      exit: '終了',
      logout: 'ログアウト'
    },
    login: {
      title: 'ログイン',
      deviceName: 'デバイス名',
      password: 'アクセスパスワード',
      submit: '入る',
      qrTip: 'モバイルでスキャンしてアクセス',
      qrAlt: 'ログインQRコード',
      deviceNameRequired: 'デバイス名を入力してください',
      failed: 'ログインに失敗しました。パスワードを確認してください'
    },
    devices: {
      columns: {
        name: 'デバイス名',
        platform: 'プラットフォーム',
        status: '状態',
        actions: '操作'
      },
      status: {
        online: 'オンライン',
        offline: 'オフライン'
      },
      actions: {
        rename: '名前変更',
        remove: '削除'
      },
      dialog: {
        renameTitle: 'デバイス名変更',
        renameInput: '新しいデバイス名を入力',
        removeTitle: 'デバイス削除',
        removeConfirm: 'このデバイスを削除してもよろしいですか？'
      },
      toast: {
        renamed: 'デバイス名を更新しました',
        removed: 'デバイスを削除しました'
      },
      entryTitle: 'アクセスURL',
      copyAddress: 'URLをコピー',
      copySuccess: 'URLをコピーしました',
      qrAlt: 'アクセスQRコード'
    },
    send: {
      title: '送信',
      selectTargets: 'デバイスを選択',
      inputText: 'テキストを入力',
      sendText: 'テキスト送信',
      selectFile: 'ファイル選択',
      sendFile: 'ファイル送信',
      toast: {
        textRequired: '送信するテキストを入力してください',
        targetRequired: '送信先デバイスを選択してください',
        fileRequired: '先にファイルを選択してください',
        textSent: 'テキストを送信しました',
        fileSent: 'ファイルを送信しました'
      }
    },
    messagesPage: {
      filters: {
        typeAll: '全タイプ',
        typeText: 'テキスト',
        typeFile: 'ファイル',
        statusAll: '全状態',
        deviceId: 'デバイスIDで絞り込み',
        sender: '送信者',
        target: '受信者',
        apply: '絞り込み',
        reset: 'リセット',
        mineOnly: '自分関連のみ'
      },
      actions: {
        batchRead: '一括既読',
        batchRetry: '一括再送',
        markRead: '既読にする',
        retry: '再送'
      },
      table: {
        time: '時間',
        type: '種類',
        content: '内容',
        download: 'ダウンロード',
        sender: '送信者',
        target: '受信者',
        status: '状態',
        operation: '操作'
      },
      status: {
        created: '未配信',
        delivered: '配信済み',
        read: '既読'
      },
      toast: {
        markReadSuccess: '既読にしました',
        markReadFailed: '既読更新に失敗しました',
        retrySuccess: '再送しました',
        retryFailed: '再送に失敗しました。後でもう一度お試しください',
        loadFailed: 'メッセージの読み込みに失敗しました',
        noneMarkable: '既読にできるメッセージがありません',
        batchReadSuccess: '{count} 件を既読にしました',
        batchReadFailed: '一括既読に失敗しました',
        selectFirst: '先にメッセージを選択してください',
        batchRetrySuccess: '{count} 件を再送しました',
        batchRetryFailed: '一括再送に失敗しました',
        batchPartial: '完了: 成功 {success} 件、失敗 {failed} 件'
      }
    },
    configPage: {
      labels: {
        port: 'サーバーポート',
        autoOpen: '自動でブラウザを開く',
        password: 'アクセスパスワード',
        language: '言語'
      },
      languages: {
        zhCN: '中国語',
        enUS: '英語',
        jaJP: '日本語'
      },
      save: '保存',
      saved: '設定を保存しました'
    },
    error: {
      AUTH_PASSWORD_INVALID: 'パスワードが正しくありません',
      AUTH_REQUIRED: '先にログインしてください',
      FILE_TOO_LARGE: 'ファイルサイズが大きすぎます',
      INVALID_REQUEST: 'リクエストパラメータが不正です',
      DEVICE_INFO_REQUIRED: 'デバイス情報が不足しています',
      DEVICE_NAME_EXISTS: 'デバイス名は既に存在します',
      DEVICE_NOT_FOUND: 'デバイスが見つかりません',
      DEVICE_NAME_REQUIRED: 'デバイス名を入力してください',
      INVALID_STATUS: '状態パラメータが不正です',
      STATUS_UPDATE_FAILED: '状態更新に失敗しました',
      INVALID_FILE: '無効なファイル形式です',
      FILE_REQUIRED: '先にファイルを選択してください',
      FILE_SAVE_FAILED: 'ファイル保存に失敗しました',
      FILE_RECORD_FAILED: 'ファイル情報の保存に失敗しました',
      FILE_NOT_FOUND: 'ファイルが見つかりません',
      MESSAGE_REQUIRED: 'メッセージ内容または送信先が必要です',
      MESSAGE_SAVE_FAILED: 'メッセージ保存に失敗しました',
      MESSAGE_LIST_FAILED: 'メッセージ一覧の取得に失敗しました',
      NOT_FOUND: '対象リソースが見つかりません',
      CONFLICT: 'リクエストが競合しています。再試行してください',
      SERVER_ERROR: 'サーバー内部エラーです。後でもう一度お試しください',
      UNKNOWN: 'リクエストに失敗しました。後でもう一度お試しください'
    }
  }
}

export default createI18n({
  legacy: false,
  locale: localeFromStorage,
  fallbackLocale: 'en-US',
  messages
})

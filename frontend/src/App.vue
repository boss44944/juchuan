<template>
  <el-container class="layout">
    <el-aside width="220px">
      <el-menu router>
        <el-menu-item index="/devices">设备</el-menu-item>
        <el-menu-item index="/messages">消息</el-menu-item>
        <el-menu-item index="/send">发送</el-menu-item>
        <el-menu-item index="/config">配置</el-menu-item>
      </el-menu>
    </el-aside>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<style scoped>
.layout {
  min-height: 100vh;
}
</style>

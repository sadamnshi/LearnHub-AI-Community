<template>
  <div id="app">
    <header >

    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: #f5f7fb;
}

#app {
  color: #0f172a;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: #ffffff;
  border-bottom: 1px solid #e2e8f0;
}

.brand {
  display: flex;
  align-items: center;
  gap: 16px;
}

.brand img {
  width: 48px;
  height: 48px;
}

.brand h1 {
  margin: 0;
  font-size: 20px;
}

.brand p {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 12px;
}

.nav {
  display: flex;
  gap: 16px;
}

.nav a {
  color: #2563eb;
  text-decoration: none;
  font-weight: 600;
}

.nav a.router-link-active {
  color: #1d4ed8;
}

main {
  padding-top: 12px;
}
</style>

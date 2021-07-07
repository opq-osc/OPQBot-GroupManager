<template>
  <div class="fullscreen bg-blue text-white text-center q-pa-md flex flex-center">
    <div>
      <div style="font-size: 15vh">
        {{ $route.query.title }}
      </div>

      <div class="text-h3" style="opacity:.4">
        {{ $route.query.info }}
      </div>

      <q-btn
        class="q-mt-xl"
        color="white"
        text-color="blue"
        unelevated
        label="联系管理员"
        no-caps
        @click="openUrl('mailto:enjoy@mcenjoy.cn')"
      />
      <q-btn
        class="q-mt-xl q-ml-sm"
        color="white"
        text-color="blue"
        unelevated
        to="/"
        label="返回首页"
        no-caps
      />
      <q-btn
        class="q-mt-xl q-ml-sm"
        color="white"
        text-color="blue"
        unelevated
        :to="$route.query.back"
        label="刷新尝试"
        no-caps
      />
    </div>
  </div>
</template>

<script>
export default {
  name: 'Error',
  methods: {
    openUrl (url) {
      window.location.href = url
    }
  }
}
</script>

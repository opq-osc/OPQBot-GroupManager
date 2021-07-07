<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="leftDrawerOpen = !leftDrawerOpen"
        />

        <q-toolbar-title>
          OPQ Bot Manager
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      :mini="miniState"
      @mouseover="miniState = false"
      @mouseout="miniState = true"
      mini-to-overlay
      show-if-above
      bordered
      content-class="bg-grey-1"
    >
      <q-list>
        <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
          :nowKey="nowKey"
        />
      </q-list>
    </q-drawer>
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
import EssentialLink from 'components/EssentialLink.vue'

const linksData = [
  {
    title: 'Github',
    caption: '欢迎Star!',
    icon: 'code',
    link: 'https://github.com/opq-osc/OPQBot-GroupManager'
  },
  {
    title: 'Mail',
    caption: 'chat.quasar.dev',
    icon: 'chat',
    link: 'mailto:enjoy@mcenjoy.cn'
  }
]

export default {
  name: 'MainLayout',
  components: { EssentialLink },
  data () {
    return {
      nowKey: 'Twitter',
      miniState: true,
      leftDrawerOpen: false,
      essentialLinks: linksData
    }
  }
}
</script>

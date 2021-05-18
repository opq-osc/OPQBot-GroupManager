<template>
  <div id="bg">
    <form autocomplete="on">
    <q-card class="login">
      <q-img height="200px"
        src="../assets/bg1.jpg"
        basic
      ><div class="absolute-bottom text-h6">
        <q-avatar v-if="avaUrl!==''" class="q-mr-md" size="lg" ><img :alt="username" v-if="avaUrl!==''" :src="avaUrl"></q-avatar>
        <q-avatar v-if="avaUrl===''" class="q-mr-md" color="primary" size="lg" text-color="white" icon="fa fa-user"></q-avatar>{{avaUrl!=='' ? title : 'OPQ 群管理机器人'  }}
      </div></q-img>
      <q-banner v-show="isError" inline-actions class="text-white bg-red">
        <q-icon size="sm" name="fa fa-exclamation-circle"></q-icon> {{ errorInfo }}
        <template v-slot:action>
          <q-btn v-show="passwordErr" flat color="white" label="忘记密码" />
          <q-btn round flat color="white" v-on:click="isError=false;passwordErr=false" icon="fa fa-times-circle"></q-btn>
        </template>
      </q-banner>
      <q-banner v-show="tips != null" inline-actions class="text-white bg-warning">
        <q-icon size="sm" name="fa fa-exclamation-circle"></q-icon> {{ tips }}
        <template v-slot:action>
          <q-btn round flat color="white" v-on:click="tips=null" icon="fa fa-times-circle"></q-btn>
        </template>
      </q-banner>
      <q-card-section>
      <q-input class="q-mb-sm" v-model="username" outlined label="用户名"></q-input>
      <q-input v-on:keypress.enter="login" autocomplete class="q-mb-sm" :type="isPwd ? 'password' : 'text'" v-model="password" outlined label="密码">
        <template v-slot:append>
        <q-icon
          :name="isPwd ? 'fa fa-eye-slash' : 'fa fa-eye'"
          class="cursor-pointer"
          @click="isPwd = !isPwd"
        />
      </template></q-input>
        <!-- <reCaptcha v-if="googleV3" class="q-mb-sm" @getValidateCode='getValidateCode'></reCaptcha> -->
      <q-toggle v-model="rememberMe" label="记住我" />
        <div class="text-right">
          <q-btn-group outline>
            <q-btn outline v-on:click="$router.replace('/')" label="返回" color="primary"/>
            <q-btn outline v-on:click="$router.replace('/register')" label="注册" color="primary"/>
            <q-btn class="float-right" v-on:click="login" label="登录" color="primary"/>
          </q-btn-group>
        </div>
      </q-card-section>

      <q-separator />
      <q-card-section>
      <a href="https://mcenjoy.cn" class="absolute-center text-caption">MCENJOY.CN</a>
      </q-card-section>
      <q-inner-loading :showing="visible">
        <q-spinner-gears size="50px" color="primary" />
      </q-inner-loading>
    </q-card></form>
  </div>
</template>

<script>
// import reCaptcha from '../layouts/google.vue'

export default {
  name: 'login',
  components: {
    // reCaptcha
  },
  data: function () {
    return {
      isPwd: true,
      rememberMe: false,
      isError: false,
      passwordErr: false,
      errorInfo: '',
      username: '',
      password: '',
      visible: true,
      avaUrl: '',
      title: '',
      token: '',
      googleV3: false,
      tokenV3: '',
      tips: null
    }
  },
  methods: {
    getValidateCode (value) {
      this.token = value
    },
    login: async function () {
      const csrf = this.$cookie('OPQWebCSRF')
      const _this = this
      if (csrf !== '') {
        if (this.username === '' || this.password === '') {
          this.isError = true
          this.errorInfo = '用户名或密码你没填呢！'
        } else {
          this.visible = true
          this.$axios.post('api/login', {
            username: this.username,
            password: this.$md5(this.password),
            csrfToken: csrf,
            rememberMe: this.rememberMe
          }).then(function (response) {
            if (response.data.code === 1) {
              _this.$store.commit('User/pushFunc', { auth: response.data.code === 1, username: response.data.data })
              _this.$q.notify({
                type: 'positive',
                position: 'top',
                message: response.data.info,
                icon: 'fa fa-check'
              })
              _this.Jump(1)
            } else if (response.data.code === 3) {
              _this.isError = true
              _this.googleV3 = true
              _this.errorInfo = response.data.info
              _this.passwordErr = false
            } else {
              _this.isError = true
              _this.errorInfo = response.data.info
              _this.passwordErr = true
            }
          }).catch(function (error) {
            _this.isError = true
            _this.errorInfo = error.message
            _this.passwordErr = false
          }).finally(function () {
            _this.visible = false
          })
        }
      } else {
        this.isError = true
        this.errorInfo = 'CSRF令牌获取失败, 无法登录!'
        this.passwordErr = false
      }
    },
    Jump: function (time) {
      setTimeout(() => {
        console.log(this.$route.query.redirect)
        if (this.$route.query.redirect !== undefined) {
          this.$router.replace(this.$route.query.redirect)
        } else {
          this.$router.replace('/')
        }
      }, time * 1000)
    }
  },
  mounted () {
    if (this.$store.state.User.auth) {
      console.log('已登录')
      this.Jump(3)
    } else {
      this.tips = this.$route.query.info
      this.visible = false
      this.$axios.get('/api/csrf')
      // this.$recaptchaLoaded().then(() => {
      //   console.log('recaptcha loaded')
      // })
    }
  }
}
</script>

<style scoped>
  #bg {
    background: url("../assets/bg.jpg");
    height: 100%;
    width: 100%;
    position: fixed;
    overflow: auto;
    background-size: cover;
  }
  .login {
    width: 80%;
    max-width: 450px;
    margin: 100px auto;
  }
</style>

<template>
  <q-page class="row">
    <div class="col-md-2 col-sm-4 col-12 q-pa-sm">
    <q-list bordered padding class="overflow-hidden rounded-borders">
      <q-item-label header>默认群配置</q-item-label>
      <q-item clickable @click="getGroupConfig(-1)">
         <q-item-section avatar>
          <q-avatar color="primary" text-color="white" icon="supervisor_account">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label>默认群配置</q-item-label>
          <q-item-label caption lines="1">当未对单个群配置时，本配置生效！</q-item-label>
        </q-item-section>
      </q-item>
       <q-item-label header>单个群管理</q-item-label>
        <q-expansion-item v-model="expanded">
           <template v-slot:header>
          <q-item-section avatar>
          <q-avatar color="primary" text-color="white">
             <img :src="'http://p.qlogo.cn/gh/' + selectGroup.GroupId + '/' + selectGroup.GroupId + '/0'">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ selectGroup.GroupName }}</q-item-label>
          <q-item-label caption lines="1">{{ selectGroup.GroupId===-1? '' : selectGroup.GroupId }}</q-item-label>
        </q-item-section>
           </template>
          <q-item  class="overflow" @click="getGroupConfig(v.GroupId)" v-for="v in groups" :key="v.GroupId" clickable v-ripple active-class="bg-grey-3" :active="v.GroupId === selectGroup.GroupId">
        <q-item-section avatar>
          <q-avatar color="primary" text-color="white">
             <img :src="'http://p.qlogo.cn/gh/' + v.GroupId + '/' + v.GroupId + '/0'">
          </q-avatar>
        </q-item-section>
        <q-item-section>
          <q-item-label>{{ v.GroupName }}</q-item-label>
          <q-item-label caption lines="1">{{ v.GroupId }}</q-item-label>
        </q-item-section>
      </q-item>
       </q-expansion-item>

    </q-list>
    </div>
    <div class="col-md-10 col-sm-8 col-12 q-pa-sm">
      <q-card flat bordered>
        <q-card-section>
        <div class="text-h6">{{ selectGroup.GroupId===-1? '默认群配置' : selectGroup.GroupName }}</div>
        <div class="text-subtitle2">{{ selectGroup.GroupId===-1? '空' :selectGroup.GroupId }}</div>
         <q-toggle class="q-mt-md" @input="setEnable" v-model="selectGroupConfig.Enable" label="启用群管理功能" />
      </q-card-section>
       <q-tabs
        v-model="tab"
        dense
        class="text-grey"
        active-color="primary"
        indicator-color="primary"
        align="left"
        inline-label
        @input="changeTab"
      >
        <q-tab name="home" icon="home" label="首页" />
        <q-tab name="members" icon="admin_panel_settings" label="群成员管理" />
        <q-tab name="entertainment" icon="mood" label="群娱乐" />
        <q-tab name="job" icon="work" label="定时任务" />
        <q-tab name="answer" icon="question_answer" label="关键词回复" />
      </q-tabs>
      <q-separator />
         <q-tab-panels v-model="tab" animated>
           <q-tab-panel name="home">
              <div class="q-col-gutter-md row items-start">
                 <q-input class="col-12 col-sm-6 col-md-4" outlined v-model="selectGroupConfig.MenuKeyWord" label="群菜单触发词 (支持正则表达式)" stack-label />
                  <q-input class="col-12 col-sm-6 col-md-4" outlined v-model="selectGroupConfig.ShutUpWord" label="群禁言触发词 (支持正则表达式)" stack-label />
                       <q-input
     outlined
     class="col-6 col-sm-3 col-md-2"
     stack-label
      v-model="selectGroupConfig.ShutUpTime"
      label="触发禁言词禁言时间(min)"
      :rules="[
        val => val >= 0 && val <= 43200 || '禁言时间最大为 43200 min',
      ]"
      lazy-rules
    >
     </q-input>
     <q-input class="col-6 col-sm-3 col-md-2" outlined v-model="selectGroupConfig.AdminUin" label="管理员QQ" stack-label />
     <q-select class="col-6 col-sm-6 col-md-2" emit-value map-options outlined v-model="selectGroupConfig.JoinVerifyType" :options="JoinVerifyTypeOptions" label="入群验证类型" />
    <q-input
     outlined
     class="col-12 col-sm-6 col-md-4"
     stack-label
      v-model="selectGroupConfig.JoinVerifyTime"
      label="入群验证码时间(s)"
      :rules="[
        val => val >= 0 || '入群图片验证码时间格式不正确',
      ]"
      lazy-rules
    ></q-input>
        <q-input
     outlined
     class="col-12 col-sm-6 col-md-4"
     stack-label
      v-model="selectGroupConfig.JoinAutoShutUpTime"
      label="入群自动禁言时间(min) 0为禁用"
      :rules="[
        val => val >= 0 && val <= 43200 || '时间格式不正确',
      ]"
      lazy-rules
    ></q-input>
                   <q-input
                    v-model="selectGroupConfig.Menu"
                    outlined
                    stack-label
                    label="菜单内容"
                    autogrow
                    filled
                    class="col-12"
                    type="textarea"
                    />
                     <q-input
                    v-model="selectGroupConfig.Welcome"
                    outlined
                    stack-label
                    label="群欢迎词 支持的宏 {name} 入群QQ名称 {uin} QQ号码"
                    autogrow
                    filled
                    class="col-12"
                    type="textarea"
                    />
              </div>
              <q-btn outline @click="setGroupConfig(selectGroup.GroupId, selectGroupConfig)" class="q-mt-md">保存</q-btn>
          </q-tab-panel>
          <q-tab-panel name="members">
             <div class="q-pa-md">
    <q-table
      grid
      title="群成员列表"
      :data="selectGroupMember.MemberList"
      :columns="columns"
      row-key="uin"
      :filter="filter"
    >
      <template v-slot:top-right>
        <q-input  outlined dense debounce="300" v-model="filter" placeholder="搜索 名称/QQ号">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>
      <template v-slot:item="props">
        <div class="q-pa-xs col-xs-12 col-sm-4 col-md-3">
          <q-card>
            <q-item>
              <q-item-section avatar>
                <q-avatar>
                  <img :src="'http://q1.qlogo.cn/g?b=qq&nk='+props.row.MemberUin+'&s=640'">
                </q-avatar>
              </q-item-section>
              <q-item-section q-item-section>
                <q-item-label>{{ props.row.NickName }} <q-badge outline :color="props.row.MemberUin === selectGroup.GroupOwner ? 'orange' : props.row.GroupAdmin === 1 ? 'green' : 'primary'" :label="props.row.MemberUin === selectGroup.GroupOwner ? '群主' : props.row.GroupAdmin === 1 ? '管理员' : '群成员'" /></q-item-label>
                <q-item-label caption>{{ props.row.MemberUin }}</q-item-label>
              </q-item-section>
            </q-item>
            <q-separator />
            <q-list>
        <q-item clickable>
          <q-item-section avatar>
            <q-icon color="primary" name="query_builder" />
          </q-item-section>

          <q-item-section>
            <q-item-label>上次发言时间</q-item-label>
            <q-item-label caption>{{ timeformat(props.row.LastSpeakTime) }}</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable>
          <q-item-section avatar>
            <q-icon color="red" name="directions_run" />
          </q-item-section>

          <q-item-section>
            <q-item-label>加入本群时间</q-item-label>
            <q-item-label caption>{{ timeformat(props.row.JoinTime) }}</q-item-label>
          </q-item-section>
        </q-item>

        <q-item clickable>
          <q-item-section avatar>
            <q-icon :color="props.row.Gender===0?'blue':props.row.Gender===1?'pink':'yellow'" name="transgender" />
          </q-item-section>

          <q-item-section>
            <q-item-label>性别</q-item-label>
            <q-item-label caption>{{ props.row.Gender===0?'小哥哥':props.row.Gender===1?'小姐姐':'未知呢？' }}</q-item-label>
          </q-item-section>
        </q-item>
      </q-list>
      <q-separator />
       <q-card-actions>
        <q-btn @click="shutup(props.row.NickName,selectGroup.GroupId,props.row.MemberUin)" outline :disable="props.row.MemberUin === selectGroup.GroupOwner ? true : props.row.GroupAdmin === 1 ? true : false">
          禁言
        </q-btn>
        <q-btn outline color="red" :disable="props.row.MemberUin === selectGroup.GroupOwner ? true : props.row.GroupAdmin === 1 ? true : false">
          移出群聊
        </q-btn>
      </q-card-actions>
          </q-card>
        </div>
      </template>
    </q-table>
  </div>
          </q-tab-panel>
          <q-tab-panel name="job">
            <div class="q-col-gutter-md row items-start">
              <div class="col-12 col-md-4">
              <q-card>
                <q-card-actions> 添加任务</q-card-actions>
                <q-separator />
                <q-card-section>
                  <div class="q-gutter-y-sm column">
                   <q-input dense outlined label="任务名称" stack-label v-model="tmp" />
                   <q-input dense outlined  label="Cron" stack-label v-model="addJob.Cron" />
                   <q-select dense emit-value map-options outlined v-model="addJob.JobType" :options="TaskOptions" label="任务类型" />

                   <q-input
                   v-model="addJob.Content"
                    outlined
                    stack-label
                    label="任务内容"
                    autogrow
                    filled
                    class="col-12"
                    type="textarea"
                    />
                  </div>
                </q-card-section>
                <q-card-actions>
                   <q-btn @click="addJobF(tmp,addJob)" outline>
                    添加
                    </q-btn>
                </q-card-actions>
              </q-card>
              </div>
               <div class="col-12 col-md-8">
  <q-card>
                <q-card-section>
                  <q-list>
      <q-item-label header>任务列表</q-item-label>

      <q-item v-for="(v, key) in selectGroupConfig.Job" :key="key">

        <q-item-section top>
          <q-item-label lines="1">
            <span class="text-weight-medium">{{ key }}</span>
            <span class="text-grey-8"> - {{ getJobTitle(v.Type) }}</span>
          </q-item-label>
          <q-item-label caption lines="1">
             {{ v.Cron }}
          </q-item-label>
        </q-item-section>

        <q-item-section top side>
          <div class="text-grey-8 q-gutter-xs">
            <q-btn size="12px" @click="delJobF(key)" color="negative" flat dense round icon="delete" />
          </div>
        </q-item-section>
      </q-item>

      <!-- <q-separator spaced /> -->

      <q-item v-show='isempty(selectGroupConfig.Job)'>

        <q-item-section top>
          <q-item-label class="q-mt-sm">任务为空</q-item-label>
        </q-item-section>
      </q-item>
    </q-list>
                </q-card-section>
              </q-card>
              </div>

            </div>
          </q-tab-panel>
<q-tab-panel name="entertainment">
 <div class="q-pa-md q-gutter-y-sm column">
   <q-toggle
      label="签到功能"
      v-model="selectGroupConfig.SignIn"
    />
     <q-toggle
      label="名片赞功能"
      v-model="selectGroupConfig.Zan"
    />
 </div>
   <q-btn outline @click="setGroupConfig(selectGroup.GroupId, selectGroupConfig)" class="q-mt-md">保存</q-btn>
    <q-btn outline class="q-mt-md q-ml-sm">导出签到数据</q-btn>
</q-tab-panel>

         </q-tab-panels>

      </q-card>
    </div>
  </q-page>
</template>

<script>
export default {
  name: 'PageIndex',
  mounted: function () {
    this.$axios.get('/api/admin/groups').then((res) => {
      if (res.data.code === 1) {
        this.groups = res.data.data.TroopList
        if (this.groups.length > 0) {
          this.getGroupConfig(-1)
        }
      } else {
        this.$q.notify({
          type: 'negative',
          position: 'top',
          message: res.data.info,
          icon: 'fa fa-check'
        })
      }
    })
  },
  data: function () {
    return {
      tmp: '',
      TaskOptions: [
        {
          label: '发送公告',
          value: 0
        },
        {
          label: '全局禁言',
          value: 1
        },
        {
          label: '全局解禁',
          value: 2
        },
        {
          label: '发送消息',
          value: 3
        }
      ],
      JoinVerifyTypeOptions: [
        {
          label: '不启用',
          value: 0
        },
        {
          label: '文字图片验证码',
          value: 1
        },
        {
          label: '算术验证码',
          value: 2
        }

      ],
      expanded: false,
      columns: [
        {
          name: 'NickName',
          required: true,
          label: '名称',
          align: 'left',
          field: 'NickName',
          format: val => `${val}`,
          sortable: true
        },
        { name: 'uin', align: 'center', label: 'QQ号', field: 'MemberUin', sortable: true }
      ],
      addJob: {},
      groups: [],
      selectGroup: {},
      selectGroupMember: {},
      selectGroupConfig: {},
      tab: 'home',
      filter: ''
    }
  },
  methods: {
    isempty: function (v) {
      let isempty = true
      for (const name in v) { // eslint-disable-line no-unused-vars
        isempty = false
      }
      return isempty
    },
    addJobF: function (k, v) {
      this.$set(this.selectGroupConfig.Job, this.tmp, this.addJob)
      this.setGroupConfig(this.selectGroup.GroupId, this.selectGroupConfig)
      this.addJob = {}
    },
    delJobF: function (k) {
      this.$q.dialog({
        title: '删除任务',
        message: '你确定要删除任务' + k + '吗?',
        cancel: '取消',
        persistent: true,
        ok: '删除'
      }).onOk(() => {
        this.$delete(this.selectGroupConfig.Job, k)
        this.setGroupConfig(this.selectGroup.GroupId, this.selectGroupConfig)
      })
    },
    getJobTitle: function (id) {
      for (let i = 0; i < this.TaskOptions.length; i++) {
        if (id === this.TaskOptions[i].value) {
          return this.TaskOptions[i].label
        }
      }
      return '未知'
    },
    shutup: function (name, id, uin) {
      this.$q.dialog({
        title: '禁言',
        message: '禁言用户' + name + '(' + uin + ') 输入禁言时间 (min)',
        prompt: {
          model: '',
          type: 'text' // optional
        },
        cancel: '取消',
        persistent: true,
        ok: '禁言'
      }).onOk(data => {
        this.$axios.post('/api/admin/shutUp', { id: id, csrfToken: this.$cookie('OPQWebCSRF'), uin: uin, time: data }).then((res) => {
          if (res.data.code === 1) {
            this.$q.notify({
              type: 'positive',
              position: 'top',
              message: res.data.info,
              icon: 'fa fa-check'
            })
          } else {
            this.$q.notify({
              type: 'negative',
              position: 'top',
              message: res.data.info,
              icon: 'fa fa-check'
            })
          }
        })
      })
    },
    getGroupMember: function (Id) {
      this.$axios.post('/api/admin/getGroupMember', { id: Id, csrfToken: this.$cookie('OPQWebCSRF') }).then((res) => {
        if (res.data.code === 1) {
          console.log(res.data)
          this.selectGroupMember = res.data.data
        } else {
          this.$q.notify({
            type: 'negative',
            position: 'top',
            message: res.data.info,
            icon: 'fa fa-check'
          })
        }
      })
    },
    changeTab: function (value) {
      if (value === 'members') {
        this.getGroupMember(this.selectGroup.GroupId)
      }
    },
    getGroupInfo: function (Id) {
      if (Id === -1) {
        this.selectGroup = { GroupId: -1 }
        this.expanded = false
        return
      }
      for (let i = 0; i < this.groups.length; i++) {
        if (this.groups[i].GroupId === Id) {
          this.selectGroup = this.groups[i]
          this.expanded = false
          return
        }
      }
    },
    timeformat: function (t) {
      return this.$timeformat(t * 1000)
    },
    setEnable: function (value) {
      this.$axios.post('/api/admin/setGroupConfig', { id: this.selectGroup.GroupId, csrfToken: this.$cookie('OPQWebCSRF'), enable: value }).then((res) => {
        if (res.data.code === 1) {
          this.selectGroup.Enable = res.data.data
        } else {
          this.$q.notify({
            type: 'negative',
            position: 'top',
            message: res.data.info,
            icon: 'fa fa-check'
          })
        }
      })
    },
    getGroupConfig: function (Id) {
      this.$axios.post('/api/admin/groupStatus', { id: Id, csrfToken: this.$cookie('OPQWebCSRF') }).then((res) => {
        if (res.data.code === 1) {
          this.getGroupInfo(Id)
          this.selectGroupConfig = res.data.data
          if (this.tab === 'members') {
            this.getGroupMember(this.selectGroup.GroupId)
          }
        } else {
          this.$q.notify({
            type: 'negative',
            position: 'top',
            message: res.data.info,
            icon: 'fa fa-check'
          })
        }
      })
    },
    setGroupConfig: function (Id, data) {
      console.log(this, data)
      this.$axios.post('/api/admin/setGroupConfig', { id: Id, csrfToken: this.$cookie('OPQWebCSRF'), data: data }).then((res) => {
        if (res.data.code === 1) {
          this.$q.notify({
            type: 'positive',
            position: 'top',
            message: res.data.info,
            icon: 'fa fa-check'
          })
        } else {
          this.$q.notify({
            type: 'negative',
            position: 'top',
            message: res.data.info,
            icon: 'fa fa-check'
          })
        }
      })
    }
  }
}
</script>

<template lang="pug">
  div.animated.fadeIn
    b-row
        b-col(md='12')
            b-card()              
              div(slot="header").pull-right                           
                b-button.btn-square(variant='success', aria-pressed='true', v-on:click="addNew") New
              b-form-fieldset()
                b-form-input(type="text" v-model="searchQuery" id="name" placeholder="Search")
              b-table.hover(:items="devices" :fields="fields")
                template(slot="edit" slot-scope="row")
                  b-button.btn-square(variant='warning' size="sm" @click.stop="editDevice(row.item)") Edit
                template(slot="del" slot-scope="row")
                  b-button.btn-square(variant='danger' size="sm" @click.stop="deleteDevice(row.item)") Delete  


</template>

<script>
import axios from 'axios'
import router from '@/router'

export default {
  name: 'devices',
  data () {
    return {
      searchQuery: '',
      devices: [],
      errors: [],
      fields: [
        {
          key: 'host', 
          sortable: true,
          label: 'Hostname',
        },
        {
          key: 'port', 
          sortable: true,
          label: 'port'         
        },
        {
          key: 'user', 
          sortable: true,
          label: 'user'         
        },
        {
          key: 'cid', 
          sortable: true,
          label: 'cid'         
        },
        {
          key: 'ws', 
          sortable: true,
          label: 'ws'         
        },
        {
          key: 'tls.enabled', 
          sortable: true,
          label: 'tls',
          // variant: 'danger'         
        },
        {
          key: 'paths[0].freq', 
          sortable: true,
          label: 'freq'         
        },
        {
          key: 'edit', 
          sortable: false,
          label: 'Edit'         
        },
         {
          key: 'del', 
          sortable: false,
          label: 'Del'         
        }          
      ]
    }
  },  
  // Fetches devices when the component is created.
  created () {
    this.getDevices()

  },
  methods: {   
    getDevices () {
      axios
      .get('http://' + window.location.hostname + ':8888/v1/devices')
      .then(response => {        
        this.devices = response.data        
      })
      .catch(err => {
        this.errors.push(err)
        console.log(err)
      })
    },
    addNew () {
      router.push('/newdevice')
    },
    deleteDevice (item) {
      console.log('aaaa', item)
      axios
        .delete('http://' + window.location.hostname + ':8888/v1/device/'+ item.uuid)
        .then(response => {        
          console.log('deleted', response.data)
          this.getDevices()
        })
        .catch(err => {
          this.errors.push(err)
          console.log(err)
        })
      
    },
    editDevice (item) {     
      router.push({
        name: 'Newdevice', 
        params: {
          item: item
        }
      })
    }
  },
  computed: {
    filteredSessions: function () {
      var self = this
      return self.devices.filter(function (session) {
        var searchRegex = new RegExp(self.searchQuery, 'i')
        return searchRegex.test(session.host) || searchRegex.test(session['peer-address'])
      })
    },
    upSessions: function () {
     
    },
    downSessions: function () {
      var self = this
      return self.devices.filter(function (session) {
        if (session['peer-state'] !== 'Established') return session
      })
    }
  }  
}
</script>

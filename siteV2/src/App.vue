<template>
  <v-app top-toolbar>
    <header>
      <v-toolbar>
        <v-btn router :to="{ name: 'Dashboard'}" icon dark>
          <v-icon>home</v-icon>
        </v-btn>
        <v-toolbar-title class="hidden-sm-and-down">terraform-ci</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items>
          <v-toolbar-item>
            <v-btn icon dark>
              <v-icon>notifications</v-icon>
            </v-btn>
          </v-toolbar-item>
          <v-toolbar-item>
            <v-btn icon dark>
              <v-icon>account_circle</v-icon>
            </v-btn>
          </v-toolbar-item>
        </v-toolbar-items>
      </v-toolbar>
    </header>
    <main>
      <v-content>
        <v-container fluid>
          <div class="ma-4">
            <router-view></router-view>
          </div>
        </v-container>
      </v-content>
    </main>
  </v-app>
</template>

<script>
  export default {
    
  }
</script>

<style lang="stylus">
  @import './stylus/theme.styl'  
</style>
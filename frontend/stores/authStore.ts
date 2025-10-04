import { defineStore } from 'pinia'
export const useAuthStore = defineStore('auth', () => {
    const loggedIn = ref(false);
    const user = ref<UserInfo | null>(null);

    function logUserIn(thisUser: UserInfo) {
        loggedIn.value = true;
        user.value = thisUser;
    }
    function logUserOut() {
        loggedIn.value = false;
        user.value = null;
    }

    return {
        loggedIn,
        user,
        logUserIn,
        logUserOut
    }
})

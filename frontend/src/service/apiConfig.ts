export const API_BASE_URL = process.env.environment == 'local' ? 'http://localhost:8080' : process.env.VUE_APP_API_BASE_URL ;
export const CHAT_URL = process.env.VUE_APP_CHAT_URL;

export default {
    API_BASE_URL,
    CHAT_URL,
};
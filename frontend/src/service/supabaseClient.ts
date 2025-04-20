import { createClient } from '@supabase/supabase-js';

// Fetch the Supabase URL and Key from environment variables
const SUPABASE_URL = process.env.VUE_APP_SUPABASE_URL!;
const SUPABASE_KEY = process.env.VUE_APP_SUPABASE_ANON_KEY!;

// Create the Supabase client
export const supabase = createClient(SUPABASE_URL, SUPABASE_KEY);

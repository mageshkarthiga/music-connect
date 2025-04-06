import { supabase } from './supabaseClient'

export interface User {
  id?: number;
  userName: string;
  emailAddress: string;
  phoneNumber: string;
  location: string;
  profilePhotoUrl: string;
}

export default {
  // Create a user
  // Example of user creation in Supabase
async createUser(user: User & { firebaseUID: string }) {
  try {
    // Sign up the user with Supabase Auth
    const { data, error } = await supabase.auth.signUp({
      email: user.emailAddress,
      password: user.phoneNumber, // Using phone number as password for simplicity
    });

    if (error) {
      throw new Error(`Error signing up: ${error.message}`);
    }

    const { user: authUser } = data;
    if (!authUser) {
      throw new Error("No user returned from sign up");
    }

    // Insert the user's details into Supabase table
    const { error: insertError } = await supabase
      .from('users')
      .insert([
        {
          user_name: user.userName,
          phone_number: user.phoneNumber,
          email_address: user.emailAddress,
          location: user.location,
          profile_photo_url: user.profilePhotoUrl,
          firebase_uid: user.firebaseUID,  // Make sure the firebaseUID is inserted correctly
        },
      ]);

    if (insertError) {
      throw new Error(`Error inserting user into the database: ${insertError.message}`);
    }

    return data; // Return the successful signup data

  } catch (error) {
    console.error('Error creating user:', error);
    throw error;
  }
},


  // Get a user by ID
  async getUser(id: number) {
    const { data, error } = await supabase
      .from('users')
      .select('*')
      .eq('id', id)
      .single();

    if (error) {
      console.error("Error fetching user:", error.message);
      throw error;
    }

    return data;
  },

  async getUserByFirebaseUID(firebaseUID: string) {
    try {
      const response = await fetch(
        `${process.env.VUE_APP_SUPABASE_URL}/rest/v1/users?select=firebase_uid,user_name,email_address,phone_number,location,profile_photo_url&firebase_uid=eq.${firebaseUID}`,
        {
          method: 'GET',
          headers: {
            'apikey': process.env.VUE_APP_SUPABASE_ANON_KEY,
            'Content-Type': 'application/json',
          },
        }
      );

      console.log("Response status:", response.status);
      console.log("Response headers:", response.headers);

      
      if (!response.ok) {
        throw new Error(`Error fetching user: ${response.statusText}`);
      }

      const data = await response.json();

      if (data.length === 0) {
        throw new Error("No user found with the provided Firebase UID");
      }

      console.log("User data:", data);
      return {
        id: data[0].id,
        userName: data[0].user_name,
        emailAddress: data[0].email_address,
        phoneNumber: data[0].phone_number,
        location: data[0].location,
        profilePhotoUrl: data[0].profile_photo_url,
      };

    } catch (error) {
      console.error("Error fetching user by Firebase UID:", error);
      throw error;
    }
  },

  // Get all users
  async getAllUsers() {
    const { data, error } = await supabase
      .from('users')
      .select('*');

    if (error) {
      console.error("Error fetching users:", error.message);
      throw error;
    }

    return data;
  },

  // Update a user
  async updateUser(id: number, updates: Partial<User>) {
    const { data, error } = await supabase
      .from('users')
      .update(updates)
      .eq('id', id);

    if (error) {
      console.error("Error updating user:", error.message);
      throw error;
    }

    return data;
  },

  // Delete a user
  async deleteUser(id: number) {
    const { error } = await supabase
      .from('users')
      .delete()
      .eq('id', id);

    if (error) {
      console.error("Error deleting user:", error.message);
      throw error;
    }
  },



  // Fetch secure data
  async fetchSecureData(idToken: string) {
    try {
      const { data, error } = await supabase
        .from('secure_data') // Assuming you have a 'secure_data' table
        .select('*')
        .eq('user_id', idToken); // Adjust the condition based on your table's structure

      if (error) {
        console.error("Error fetching secure data:", error.message);
        throw error;
      }

      return data;
    } catch (error) {
      console.error("Error fetching secure data:", error);
      throw error;
    }
  },
};

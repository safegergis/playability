<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-1/2 min-w-96 shadow-md rounded-md p-6">
      <CardHeader>
        <CardTitle class="text-2xl font-semibold text-center">Login</CardTitle>
      </CardHeader>
      <CardContent>
        <Form v-slot="{ meta }" :validation-schema="schema" @submit="onSubmit">
          <div class="mb-4">
            <div class="mb-4">
              <label for="email" class="block text-gray-700 mb-2">Email</label>
              <Field
                id="email"
                name="email"
                type="email"
                placeholder="Enter your email"
                class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                @input="InvalidLogin = false"
              />
              <ErrorMessage name="email" class="text-red-500 text-sm mt-1" />
            </div>

            <div class="mb-4">
              <label for="password" class="block text-gray-700 mb-2"
                >Password</label
              >
              <Field
                id="password"
                name="password"
                type="password"
                placeholder="Enter your password"
                class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
                @input="InvalidLogin = false"
              />
              <ErrorMessage name="password" class="text-red-500 text-sm mt-1" />
            </div>
            <p v-if="InvalidLogin" class="text-red-500 text-sm mt-1">
              Invalid email or password
            </p>
            <Button type="submit" class="w-full" :disabled="!meta.valid">
              Register
            </Button>
          </div></Form
        >
      </CardContent>
    </Card>
    <Toaster />
  </div>
</template>

<script lang="ts" setup>
import { Form, Field, ErrorMessage } from "vee-validate";
import * as yup from "yup";
import { useToast } from "@/components/ui/toast";
const { toast } = useToast();

const InvalidLogin = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
  email: yup
    .string()
    .required("Email is required")
    .email("Must be a valid email"),
  password: yup
    .string()
    .required("Password is required")
    .min(6, "Password must be at least 6 characters"),
});

// Handle form submission
const onSubmit = async (values: Record<string, string>) => {
  const { data, error } = await useAsyncData("login", () =>
    $fetch("/api/auth/login", {
      method: "POST",
      body: values,
    })
  );
  // Check if the error is a 401
  if (error.value && error.value.statusCode !== 401) {
    //if not 401 then show error
    console.log(error.value);
    console.error("Login error:", error.value);
    toast({
      title: "Login failed!",
      description: "Please try again.",
    });
  } else if (error.value && error.value.statusCode === 401) {
    //if 401 then show invalid login
    console.log("error.value.message", error.value.data);
    InvalidLogin.value = true;
  } else {
    //success
    console.log("Login successful! ", data.value);
    toast({
      title: "Login successful!",
      description: "You have successfully logged in.",
    });
    ///navigateTo("/");
  }
};
</script>

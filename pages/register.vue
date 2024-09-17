<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-1/2 min-w-96 shadow-md rounded-md p-6">
      <CardHeader>
        <CardTitle class="text-2xl font-semibold text-center"
          >User Registration</CardTitle
        >
      </CardHeader>
      <CardContent>
        <Form v-slot="{ meta }" :validation-schema="schema" @submit="onSubmit">
          <div class="mb-4">
            <label for="username" class="block text-gray-700 mb-2"
              >Username</label
            >
            <p class="text-gray-500 text-sm mb">
              This is your public display name
            </p>
            <Field
              id="username"
              name="username"
              type="text"
              placeholder="Enter your username"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
              @input="uniqueUsername = false"
            />
            <ErrorMessage name="username" class="text-red-500 text-sm mt-1" />
            <p v-if="uniqueUsername" class="text-red-500 text-sm mt-1">
              Username is already in use
            </p>
          </div>

          <div class="mb-4">
            <label for="email" class="block text-gray-700 mb-2">Email</label>
            <Field
              id="email"
              name="email"
              type="email"
              placeholder="Enter your email"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
              @input="uniqueEmail = false"
            />
            <ErrorMessage name="email" class="text-red-500 text-sm mt-1" />
            <p v-if="uniqueEmail" class="text-red-500 text-sm mt-1">
              Email is already in use
            </p>
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
            />
            <ErrorMessage name="password" class="text-red-500 text-sm mt-1" />
          </div>

          <div class="mb-6">
            <label for="confirmPassword" class="block text-gray-700 mb-2"
              >Confirm Password</label
            >
            <Field
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              placeholder="Confirm your password"
              class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-400"
            />
            <ErrorMessage
              name="confirmPassword"
              class="text-red-500 text-sm mt-1"
            />
          </div>

          <Button type="submit" class="w-full" :disabled="!meta.valid">
            Register
          </Button>
        </Form>
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

const uniqueEmail = ref(false);
const uniqueUsername = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
  username: yup
    .string()
    .required("Username is required")
    .min(3, "Username must be at least 3 characters"),
  email: yup
    .string()
    .required("Email is required")
    .email("Must be a valid email"),
  password: yup
    .string()
    .required("Password is required")
    .min(6, "Password must be at least 6 characters"),
  confirmPassword: yup
    .string()
    .required("Please confirm your password")
    .oneOf([yup.ref("password")], "Passwords must match"),
});

// Handle form submission
const onSubmit = async (values: {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}) => {
  const registrationData = {
    username: values.username,
    email: values.email,
    password: values.password,
  };
  // Reset form values after submission

  const { error } = await useAsyncData("register", () =>
    $fetch("http://localhost:8080/user/register", {
      method: "POST",
      body: registrationData,
    })
  );
  if (error.value && error.value.statusCode !== 409) {
    console.log(error.value);
    console.error("Registration error:", error.value);
    toast({
      title: "Registration failed!",
      description: "Please try again.",
    });
  } else if (error.value && error.value.statusCode === 409) {
    console.log("error.value.message", error.value.data);
    const cleanMessage = (error.value.data as string).trim();
    if (cleanMessage === "email is already in use") {
      uniqueEmail.value = true;
    } else if (cleanMessage === "username is already in use") {
      uniqueUsername.value = true;
    }
  } else {
    toast({
      title: "Registration successful!",
      description: "You have successfully registered.",
    });
    ///navigateTo("/");
  }
};
</script>

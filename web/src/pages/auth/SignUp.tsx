import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { gql, useMutation } from '@apollo/client';
import { Link, useNavigate } from 'react-router-dom';
import type { SignUpInput } from '../../types';
import { PROJECT_NAME } from '../../lib/utils';

const SIGN_UP = gql`
  mutation SignUp($input: AuthSignUpByPasswordInput!) {
    signUpByPassword(input: $input) {
      success
      accessToken
    }
  }
`;

const quotes = [
  "The real voyage of discovery consists not in seeking new landscapes, but in having new eyes.",
  "Discovery consists of seeing what everybody has seen and thinking what nobody has thought.",
  "The greatest discovery of all time is that a person can change his future by merely changing his attitude.",
  "The important thing is not to stop questioning. Curiosity has its own reason for existing.",
  "Every great advance in science has issued from a new audacity of imagination.",
  "The more I read, the more I acquire, the more certain I am that I know nothing.",
  "The only true wisdom is in knowing you know nothing.",
  "The discovery of a new dish does more for the happiness of mankind than the discovery of a star.",
  "We don't receive wisdom; we must discover it for ourselves after a journey that no one can take for us.",
  "The most beautiful thing we can experience is the mysterious. It is the source of all true art and science.",
  "To see the world, things dangerous to come to, to see behind walls, draw closer, to find each other, and to feel. That is the purpose of life.",
  "Man cannot discover new oceans unless he has the courage to lose sight of the shore.",
  "The greatest discoveries have come from people who have looked at a standard situation and seen it differently.",
  "The best way to predict the future is to invent it.",
  "The universe is full of magical things patiently waiting for our wits to grow sharper.",
  "The greatest obstacle to discovery is not ignorance - it is the illusion of knowledge.",
  "In wisdom gathered over time I have found that every experience is a form of exploration.",
  "Not all those who wander are lost.",
  "The more I learn, the more I realize how much I don't know.",
  "Exploration is really the essence of the human spirit."
];

export function SignUp() {
  const navigate = useNavigate();
  const [quote, setQuote] = useState("");
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<SignUpInput>();
  const [signUp] = useMutation(SIGN_UP);

  useEffect(() => {
    const randomIndex = Math.floor(Math.random() * quotes.length);
    setQuote(quotes[randomIndex]);
  }, []);

  const onSubmit = async (data: SignUpInput) => {
    try {
      const result = await signUp({
        variables: { input: data }
      });
      
      if (result.data?.signUpByPassword.success) {
        const token = result.data.signUpByPassword.accessToken;
        localStorage.setItem('authToken', token);
        navigate('/dashboard');
      }
    } catch (error) {
      console.error('Sign up error:', error);
    }
  };

  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      <div className="w-full md:w-1/2 bg-gradient-to-b from-green-500 to-teal-600 animate-gradient-y flex items-end justify-center">
        <div className="p-4">
          <p className="text-white text-lg font-semibold">{quote}</p>
        </div>
      </div>
      <div className="w-full md:w-1/2 flex flex-col justify-center items-center bg-white p-8">
        <div className="max-w-sm w-full space-y-6">
          <h2 className="text-3xl font-bold text-center text-gray-900">Join <span className="bg-gradient-to-r from-blue-500 to-purple-600 text-transparent bg-clip-text animate-pulse">{PROJECT_NAME}</span></h2>
          <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
            <div>
              <label htmlFor="firstName" className="sr-only">First name</label>
              <input
                {...register('firstName', { required: 'First name is required' })}
                type="text"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="First name"
              />
              {errors.firstName && <p className="text-red-500 text-xs mt-1">{errors.firstName.message}</p>}
            </div>
            <div>
              <label htmlFor="lastName" className="sr-only">Last name</label>
              <input
                {...register('lastName', { required: 'Last name is required' })}
                type="text"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Last name"
              />
              {errors.lastName && <p className="text-red-500 text-xs mt-1">{errors.lastName.message}</p>}
            </div>
            <div>
              <label htmlFor="email" className="sr-only">Email address</label>
              <input
                {...register('email', { required: 'Email is required' })}
                type="email"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Email address"
              />
              {errors.email && <p className="text-red-500 text-xs mt-1">{errors.email.message}</p>}
            </div>
            <div>
              <label htmlFor="password" className="sr-only">Password</label>
              <input
                {...register('password', { required: 'Password is required' })}
                type="password"
                className="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="Password"
              />
              {errors.password && <p className="text-red-500 text-xs mt-1">{errors.password.message}</p>}
            </div>
            <div>
              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                Sign Up
              </button>
            </div>
          </form>
          <div className="text-center">
            <Link to="/signin" className="text-indigo-600 hover:text-indigo-500">Already have an account? Sign in</Link>
          </div>
        </div>
      </div>
      <style jsx>{`
        @keyframes transform {
          0%, 100% {
            transform: translateY(0);
          }
          50% {
            transform: translateY(-10px);
          }
        }
        .animate-transform {
          animation: transform 3s ease-in-out infinite;
        }
        @keyframes expressive {
          0%, 100% {
            transform: translateY(0) scale(1);
          }
          25% {
            transform: translateY(-20px) scale(1.1);
          }
          50% {
            transform: translateY(10px) scale(0.9);
          }
          75% {
            transform: translateY(-10px) scale(1.05);
          }
        }
        .animate-expressive {
          animation: expressive 5s ease-in-out infinite;
        }
      `}</style>
    </div>
  );
}
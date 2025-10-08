import 'package:flutter/material.dart';

class GuideItem {
  final String? imagePath;
  final String? videoAsset;
  final bool hasVideo;
  final String title;
  final String description;

  const GuideItem({
    required this.title,
    required this.description,
    this.imagePath,
    this.videoAsset,
    this.hasVideo = false,
  });
}

class GuideModule {
  final int id;
  final String title;
  final String description;
  final String imagePath;
  final List<GuideItem> items;

  const GuideModule({
    required this.id,
    required this.title,
    required this.description,
    required this.imagePath,
    required this.items,
  });
}

final List<GuideModule> allGuideModules = [
  GuideModule(
    id: 0,
    title: 'Login and Registration',
    description: 'Understand how to sign in, create an account, and verify your email to unlock LogiQ.',
    imagePath: 'assets/guide/images/Login and Registration.png',
    items: [
      GuideItem(
        title: 'Login',
        description: 'Existing users can log in by entering their email and password.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Login.mp4',
      ),
      GuideItem(
        title: 'Register',
        description: 'Create a new account with your email, username, and a secure password that meets the requirements.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Register.mp4',
      ),
      GuideItem(
        title: 'Verification',
        description: 'After registration, enter the verification code sent via email to complete the setup.',
        hasVideo: false,
        imagePath: 'assets/guide/images/Verification.png',
      ),
    ],
  ),
  GuideModule(
    id: 1,
    title: 'Home',
    description: 'Start quizzes, monitor your overall accuracy, and switch themes from the home page.',
    imagePath: 'assets/guide/images/Home.png',
    items: [
      GuideItem(
        title: 'Start a Quiz',
        description: 'Choose from quizzes that target different knowledge points and difficulty levels, then tap one to begin.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Start-a-Quiz.mp4',
      ),
      GuideItem(
        title: 'Performance Statistics',
        description: 'Review your average daily accuracy at a glance, with deeper insights available on the Profile page.',
        hasVideo: false,
        imagePath: 'assets/guide/images/Performance Statistics.png',
      ),
      GuideItem(
        title: 'Offline Mode',
        description: 'Answer quizzes without an internet connection using the local question bank; results are not recorded or analyzed.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Offline-Mode.mp4',
      ),
      GuideItem(
        title: 'Dark Mode Switch',
        description: 'Toggle between light and dark themes from the top-right corner for a comfortable visual experience.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Dark-Mode-Switch.mp4'
      ),
    ],
  ),
  GuideModule(
    id: 2,
    title: 'Quiz Features',
    description: 'Learn how quizzes track progress, cover multiple knowledge dimensions, and deliver results.',
    imagePath: 'assets/guide/images/Quiz Features.png',
    items: [
      GuideItem(
        title: 'Timer and Progress Tracking',
        description: 'Keep an eye on the countdown timer and progress bar at the top of the Quiz page to manage your pace.',
        hasVideo: false,
        imagePath: 'assets/guide/images/Timer and Progress Tracking.png',
      ),
      GuideItem(
        title: 'Three Dimensions',
        description: 'Work across Truth Table, Equivalence, and Inference dimensions, three difficulty levels, and multiple question types.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Three-Dimensions.mp4',
      ),
      GuideItem(
        title: 'Submission and Results',
        description: 'Submit once all questions are answered to view accuracy and a detailed performance breakdown.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Submission-and-Results.mp4',
      ),
    ],
  ),
  GuideModule(
    id: 3,
    title: 'Review',
    description: 'Revisit completed quizzes and track your learning journey with detailed history.',
    imagePath: 'assets/guide/images/Review.png',
    items: [
      GuideItem(
        title: 'Single Quiz Review',
        description: 'Inspect the time spent, accuracy, and correct answers for each completed quiz.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Single-Quiz-Review.mp4',
      ),
      GuideItem(
        title: 'Quiz History',
        description: 'Browse all past quizzes in History to monitor progress over time.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Quiz-History.mp4',
      ),
    ],
  ),
  GuideModule(
    id: 4,
    title: 'Profile',
    description: 'Manage personal details, security settings, and a detailed overview of your performance.',
    imagePath: 'assets/guide/images/Profile.png',
    items: [
      GuideItem(
        title: 'Statistics Overview',
        description: 'Review daily accuracy, error distribution, completed quizzes, and time spent in one place.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Statistics-Overview.mp4',
      ),
      GuideItem(
        title: 'Edit Profile',
        description: 'Update your avatar, username, and other personal information whenever needed.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Edit-Profile.mp4',
      ),
      GuideItem(
        title: 'Change Password',
        description: 'Secure your account by updating the password, including email verification for added safety.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Change-Password.mp4',
      ),
      GuideItem(
        title: 'Logout',
        description: 'Sign out safely and re-enter your credentials when returning to the app.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Logout.mp4',
      ),
    ],
  ),
  GuideModule(
    id: 5,
    title: 'Admin Features',
    description: 'Administrative tools for managing the question bank and analytics (admin access only).',
    imagePath: 'assets/guide/images/Admin Features.png',
    items: [
      GuideItem(
        title: 'Question Bank Management',
        description: 'View the full question list with accuracy stats, apply filters, and handle pagination as needed.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Question-Bank-Management.mp4',
      ),
      GuideItem(
        title: 'Delete or Export Questions',
        description: 'Remove outdated questions or export selected ones for external review.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Delete-or-Export-Questions.mp4',
      ),
      GuideItem(
        title: 'Question Generation',
        description: 'Configure quantity and difficulty to auto-generate questions, then preview or export the results.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Question-Generation.mp4',
      ),
      GuideItem(
        title: 'Statistics',
        description: 'Use dashboards that visualize accuracy and category distribution to refine the question bank.',
        hasVideo: true,
        videoAsset: 'assets/guide/videos/Statistics.mp4',
      ),
    ],
  ),
];

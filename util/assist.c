#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <curl/curl.h>

#define MAX_MOVIES 50000 

typedef struct {
  char titleType[50];
  char primaryTitle[200];
  int startYear;
  char runtimeMinutes[10];
  char genres[100];
  float averageRating;
  int numVotes;
} Movie;


void printMovie(const Movie* movie) {
    printf("Title: %s\n", movie->primaryTitle);
    printf("Type: %s\n", movie->titleType);
    printf("Year: %d\n", movie->startYear);
    printf("Runtime: %s mins\n", movie->runtimeMinutes);
    printf("Genres: %s\n", movie->genres);
    printf("Rating: %.2f\n", movie->averageRating);
    printf("Votes: %d\n", movie->numVotes);
    printf("--------------------------\n");
}
static size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp)
{
    (void)contents;  // Prevent unused variable warning
    (void)userp;     // Prevent unused variable warning
    return size * nmemb;
}

void send_post_request(const Movie* movie) {
  CURL* curl;
  CURLcode res;

  curl_global_init(CURL_GLOBAL_DEFAULT);
  curl = curl_easy_init();

  if(curl) {
    struct curl_slist* headers = NULL;

    // Set the URL
    curl_easy_setopt(curl, CURLOPT_URL, "http://localhost:8080/addmovie");
    
    // Add headers
    headers = curl_slist_append(headers, "Content-Type: application/json");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

    // Construct the POST data from the movie struct
    char post_data[512]; // Make sure this is large enough for your movie data
    snprintf(post_data, sizeof(post_data), 
          "{\"titleType\": \"%s\", \"primaryTitle\": \"%s\", \"startYear\": \"%d\", \"runTimeMinutes\": \"%s\", \"genres\": \"%s\", \"averageRating\": \"%.1f\", \"numVotes\": \"%d\"}",
          movie->titleType, movie->primaryTitle, movie->startYear, movie->runtimeMinutes, movie->genres, movie->averageRating, movie->numVotes);

    
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, post_data);

    // Some settings to ensure we get the response
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);

    // Perform the request
    res = curl_easy_perform(curl);
    if(res != CURLE_OK) {
      fprintf(stderr, "curl_easy_perform() failed: %s\n", curl_easy_strerror(res));
    }

      // Cleanup
    curl_easy_cleanup(curl);
    curl_slist_free_all(headers);
  }

  curl_global_cleanup();
}

int main() {
  clock_t file_start, file_end, req_start, req_end;
  double file_cpu_time_used, req_cpu_time_used;
  file_start = clock();
  
  FILE* movie_info;
  movie_info = fopen("output.tsv", "r");


  if (movie_info == NULL) {
    printf("file cannot be properly opened. Exiting \n");
    exit(0);
  }  

  // Skip the header line
  char header[1024];  // Assuming the header won't exceed 1000 characters
  fgets(header, sizeof(header), movie_info);
  Movie* movies = (Movie*) malloc(MAX_MOVIES * sizeof(Movie));
  if(movies == NULL) {
    printf("Something went wrong with allocating memory for the movie list \n "); 
    exit(0);
  }
  int movie_count=0;
  while (movie_count < MAX_MOVIES &&
  fscanf(movie_info,"%*s\t%49[^\t]\t%199[^\t]\t%*199[^\t]\t%*d\t%d\t%*s\t%9[^\t]\t%99[^\t]\t%f\t%d", 
            movies[movie_count].titleType, 
            movies[movie_count].primaryTitle, 
            &movies[movie_count].startYear, 
            movies[movie_count].runtimeMinutes, 
            movies[movie_count].genres, 
            &movies[movie_count].averageRating, 
            &movies[movie_count].numVotes) > 1 ) { 
              // Ensure all 7 fields were read your loop's body
              // printf("Movie %d ", movie_count);
              // printMovie(&movies[movie_count]); 
              movie_count++;
            }

  file_end = clock();
  file_cpu_time_used = ((double) (file_end - file_start)) / CLOCKS_PER_SEC; 
  //do some request stuff
  req_start = clock();
  int movies_within_criteria = 0;
  for(int i=0; i< MAX_MOVIES; i++) {
    //i want all movies > 6.8 rating and english? and Movie and > 1970 
    if(strcmp(movies[i].titleType, "movie") != 0){
      //printf("not a movie : %s", movies[i].titleType);
      continue;
    } 
    if(movies[i].startYear < 1970) {
      //printf("not a good year %d", movies[i].startYear); 
      continue;
    }
    if(movies[i].averageRating < 5.0) {
      //printf("not a good rating %.1f", movies[i].averageRating); 
      continue;
    }
    if(movies[i].numVotes < 1000) {
      continue;
    }
    movies_within_criteria++; 
    //printMovie(&movies[i]); 
    send_post_request(&movies[i]);
  }
  req_end = clock(); 
  req_cpu_time_used = ((double) (req_end - req_start)) / CLOCKS_PER_SEC;
  printf("Timing - File time used: %.2fms - Request time used: %.2fms \n", file_cpu_time_used * 1000, req_cpu_time_used * 1000);
  printf("freeing \n");
  free(movies); 
  
  return 0;
}

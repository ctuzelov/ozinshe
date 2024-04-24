package helper

import (
	"fmt"
	"ozinshe/internal/models"
)

func HasMatchingGenres(movieGenres []string, filterGenres map[string]struct{}) bool {
	for _, movieGenre := range movieGenres {
		if _, ok := filterGenres[movieGenre]; ok {
			return true
		}
	}
	return false
}

func GetMovieGenres(movie models.Movie) []string {
	genres := []string{}
	for _, genre := range movie.Genres {
		genres = append(genres, genre.Name)
	}
	return genres
}

func GetSeriesGenres(series models.Series) []string {
	genres := []string{}
	for _, genre := range series.Genres {
		genres = append(genres, genre.Name)
	}
	return genres
}

func FilterAndRemoveMoviesDuplicates(movies []models.Movie, filter models.FilterParams) ([]models.Movie, error) {
	filteredMovies := []models.Movie{}

	genresMap := map[string]struct{}{}
	for _, genre := range filter.Genres {
		genresMap[genre] = struct{}{}
	}

	if len(filter.Genres) > 0 {
		for _, movie := range movies {
			movieGenres := GetMovieGenres(movie)

			if HasMatchingGenres(movieGenres, genresMap) {
				filteredMovies = append(filteredMovies, movie)
			}
		}
	} else {
		filteredMovies = movies
	}

	fmt.Println(" Filtered Movies line 53 ", filteredMovies)
	if filter.YearStart != 0 || filter.YearEnd != 0 {
		var tempMovies []models.Movie
		for _, movie := range filteredMovies {
			if filter.YearStart <= movie.ReleaseYear && movie.ReleaseYear <= filter.YearEnd {
				tempMovies = append(tempMovies, movie)
			}
		}
		filteredMovies = append(filteredMovies, tempMovies...)
	}

	fmt.Println(" Filtered Movies line 62 ", filteredMovies)
	seen := make(map[int]struct{})
	var uniqueMovies []models.Movie
	for _, movie := range filteredMovies {
		if _, ok := seen[movie.ID]; !ok {
			uniqueMovies = append(uniqueMovies, movie)
			seen[movie.ID] = struct{}{}
		}
	}

	fmt.Println(" Unique Movies ", uniqueMovies)
	return uniqueMovies, nil
}

func FilterAndRemoveSeriesDuplicates(series []models.Series, filter models.FilterParams) ([]models.Series, error) {
	filteredSeries := []models.Series{}

	genresMap := map[string]struct{}{}
	for _, genre := range filter.Genres {
		genresMap[genre] = struct{}{}
	}

	// Genre Filtering
	if len(filter.Genres) > 0 {
		for _, series := range series {
			seriesGenres := GetSeriesGenres(series)

			if HasMatchingGenres(seriesGenres, genresMap) {
				filteredSeries = append(filteredSeries, series)
			}
		}
	} else {
		filteredSeries = series
	}

	if filter.YearStart != 0 || filter.YearEnd != 0 {
		var tempSeries []models.Series
		for _, series := range filteredSeries {
			if filter.YearStart <= series.ReleaseYear && series.ReleaseYear <= filter.YearEnd {
				tempSeries = append(tempSeries, series)
			}
		}
		filteredSeries = append(filteredSeries, tempSeries...)
	}

	seen := make(map[int]struct{})
	var uniqueSeries []models.Series
	for _, series := range filteredSeries {
		if _, ok := seen[series.ID]; !ok {
			uniqueSeries = append(uniqueSeries, series)
			seen[series.ID] = struct{}{}
		}
	}

	return uniqueSeries, nil
}

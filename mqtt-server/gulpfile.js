var gulp = require('gulp')
var babel = require('gulp-babel')

gulp.task('default', function () {
  return gulp.src(['**/*.js', '!dist/**', '!node_modules/**', '!gulpfile.js'])
    .pipe(babel())
    .pipe(gulp.dest('dist'))
})

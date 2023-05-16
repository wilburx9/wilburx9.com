const {series, watch, src, dest, parallel} = require('gulp');
const pump = require('pump');

const GhostAdminAPI = require("@tryghost/admin-api");

// gulp plugins and utils
var livereload = require('gulp-livereload');
var postcss = require('gulp-postcss');
var zip = require('gulp-zip');
var uglify = require('gulp-uglify');
var beeper = require('beeper');

// postcss plugins
var autoprefixer = require('autoprefixer');
var colorFunction = require('postcss-color-mod-function');
var cssnano = require('cssnano');
var easyimport = require('postcss-easy-import');
var tailwindcss = require('tailwindcss');

function serve(done) {
    livereload.listen();
    done();
}

const handleError = (done) => {
    return function (err) {
        if (err) {
            beeper();
        }
        return done(err);
    };
};

function hbs(done) {
    pump([
        src(['*.hbs', '**/**/*.hbs', '!node_modules/**/*.hbs']),
        livereload()
    ], handleError(done));
    css(done);
}

function css(done) {
    var processors = [
        easyimport,
        colorFunction(),
        tailwindcss(),
        autoprefixer(),
        cssnano()
    ];

    pump([
        src('assets/css/*.css', {sourcemaps: true}),
        postcss(processors),
        dest('assets/built/', {sourcemaps: '.'}),
        livereload()
    ], handleError(done));
}

function js(done) {
    pump([
        src('assets/js/*.js', {sourcemaps: true}),
        uglify(),
        dest('assets/built/', {sourcemaps: '.'}),
        livereload()
    ], handleError(done));
}

function zipper(done) {
    var targetDir = 'dist/';
    var themeName = require('./package.json').name;
    var filename = themeName + '.zip';

    pump([
        src([
            '**',
            '!node_modules',
            '!node_modules/**',
            '!dist',
            '!dist/**',
            '!**/*.map',
            '!assets/css/**',
            '!assets/js/**',
            '!assets/screenshot-desktop.jpg'
        ]),
        zip(filename),
        dest(targetDir)
    ], handleError(done));
}

async function deploy(done) {
    let url = process.env.GHOST_API_URL
    let apiKey = process.env.GHOST_API_KEY
    let themeName = require('./package.json').name
    let zipPath = `dist/${themeName}.zip`
    let admin = new GhostAdminAPI({
        url: url,
        key: apiKey,
        version: "v5"
    })
    await admin.themes.upload({file: zipPath})
    await admin.themes.activate(themeName)
    done()
}

const cssWatcher = () => watch('assets/css/**', css);
const jsWatcher = () => watch('assets/js/**', js);
const hbsWatcher = () => watch(['*.hbs', '**/**/*.hbs', '!node_modules/**/*.hbs'], hbs);
const watcher = parallel(cssWatcher, hbsWatcher, jsWatcher);
const build = series(css, js);
const dev = series(build, serve, watcher);

exports.build = build;
exports.zip = series(build, zipper);
exports.default = dev;
exports.deploy = deploy

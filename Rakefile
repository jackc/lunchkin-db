begin
  require 'bundler'
  Bundler.setup
rescue LoadError
  puts 'You must `gem install bundler` and `bundle install` to run rake tasks'
end

require 'fileutils'
require 'rspec/core/rake_task'

file 'lunchkin-db' => ['handlers.go', 'main.go', 'migrations.go'] do |t|
  sh 'go build'
end

desc 'Build lunchkin-db'
task build: 'lunchkin-db'

desc 'Run lunchkin-db server'
task server: 'lunchkin-db' do
  sh './lunchkin-db'
end

task spec_server: :build do
  FileUtils.mkdir_p 'tmp/spec/server'
  FileUtils.touch 'tmp/spec/server/stdout.log'
  FileUtils.touch 'tmp/spec/server/stderr.log'
  pid = Process.spawn './lunchkin-db -config=config.test.yml',
    out: 'tmp/spec/server/stdout.log',
    err: 'tmp/spec/server/stderr.log'
  at_exit { Process.kill 'TERM', pid }
end

RSpec::Core::RakeTask.new(:spec)
task spec: :spec_server

task :default => :spec

begin
  require 'bundler'
  Bundler.setup
rescue LoadError
  puts 'You must `gem install bundler` and `bundle install` to run rake tasks'
end

require 'rspec/core/rake_task'

file 'lunchkin-db' => ['handlers.go', 'main.go', 'migrations.go'] do |t|
  sh 'go build'
end

desc 'Build lunchkin-db'
task build: 'lunchkin-db'

RSpec::Core::RakeTask.new(:spec)
task spec: :build

task :default => :spec

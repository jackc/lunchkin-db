require 'spec_helper'

feature 'Game entry' do
  def enter_player(player_name, level, effective_level, winner)
    check player_name

    level_input = find :xpath, ".//td/label[text()='#{player_name}']/ancestor::tr//input[@ng-model='p.level']"
    level_input.set level

    effective_level_input = find :xpath, ".//td/label[text()='#{player_name}']/ancestor::tr//input[@ng-model='p.effective_level']"
    effective_level_input.set effective_level

    winner_input = find :xpath, ".//td/label[text()='#{player_name}']/ancestor::tr//input[@ng-model='p.winner']"
    winner_input.set winner
  end



  background  do
    DB[:player].insert name: 'Frodo'
    DB[:player].insert name: 'Samwise'
    DB[:player].insert name: 'Gandolf'
  end

  scenario 'Entering a game' do
    click_on 'Record a Game'

    fill_in 'Date', with: '2013-08-19'
    fill_in 'Length', with: '5'

    enter_player('Frodo', 8, 15, false)
    enter_player('Samwise', 9, 20, false)
    enter_player('Gandolf', 10, 25, true)

    click_on 'Save'

    expect(DB[:game].count).to eq 1
    expect(DB[:game_player].count).to eq 3
  end
end

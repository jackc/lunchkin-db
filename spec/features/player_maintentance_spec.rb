require 'spec_helper'

feature 'Player maintentance' do
  scenario 'Adding a player' do
    click_on 'Players'

    fill_in 'player_name', with: 'Frodo'
    click_on 'Add'

    within '#playersPage ul' do
      expect(page).to have_content 'Frodo'
    end
    expect(DB[:player].count).to eq 1
  end

  scenario 'Deleting a player' do
    DB[:player].insert name: 'Frodo'
    click_on 'Players'

    within '#playersPage ul' do
      delete_button = find :xpath, ".//li/div[contains(@class, 'name')]/../button"
      delete_button.click
    end

    expect(page).to_not have_content 'Frodo'
    expect(DB[:player].count).to eq 0
  end
end
